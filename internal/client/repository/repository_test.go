package repository_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/internal/client/storage"
	"github.com/niksmo/gophkeeper/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type suite struct {
	ctx context.Context
	r   *repository.Repository
	s   *storage.Storage
}

type repoConstructor func(logger.Logger, repository.Storage) *repository.Repository

func newSuite(t *testing.T, newRepoFn repoConstructor) *suite {
	log := logger.NewPretty("debug")
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()

	t.Cleanup(func() {
		os.Remove(dsn)
	})

	s := storage.New(log, dsn)
	s.MustRun(ctx)
	r := newRepoFn(log, s)
	return &suite{ctx, r, s}
}

func TestAdd(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		expectedID := 1
		expectedName := "testName"
		expectedData := []byte("helloWorld")
		id, err := st.r.Create(st.ctx, expectedName, expectedData)
		require.NoError(t, err)
		assert.Equal(t, expectedID, id)

		stmt := `
		SELECT name, data FROM passwords;
		`
		rows, err := st.s.QueryContext(st.ctx, stmt)
		require.NoError(t, err)
		defer rows.Close()
		var nRows int
		var name string
		var data []byte
		for rows.Next() {
			nRows++
			err := rows.Scan(&name, &data)
			require.NoError(t, err)
		}
		err = rows.Err()
		require.NoError(t, err)
		require.Equal(t, 1, nRows)
		assert.Equal(t, expectedName, name)
		assert.Equal(t, expectedData, data)
	})

	t.Run("UniqueNameConstraintErr", func(t *testing.T) {
		st := newSuite(t, repository.NewBin)
		objectName := "testName"
		objectData := []byte("testData")
		expectedID := 1
		actualID, err := st.r.Create(st.ctx, objectName, objectData)
		require.NoError(t, err)
		require.Equal(t, expectedID, actualID)

		_, err = st.r.Create(st.ctx, objectName, objectData)
		assert.Error(t, err)
	})
}

func TestReadByID(t *testing.T) {
	t.Run("ExistsID", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		stmt := `
		INSERT INTO
		  passwords (name, data, created_at, updated_at)
		VALUES (?, ?, ?, ?);
		`
		name := "testName"
		data := []byte("helloTestWorld")
		createdAt := time.Now()
		updatedAt := createdAt
		_, err := st.s.ExecContext(
			st.ctx, stmt, name, data, createdAt, updatedAt,
		)
		require.NoError(t, err)

		id := 1
		retData, err := st.r.ReadByID(st.ctx, id)
		require.NoError(t, err)
		assert.Equal(t, data, retData)
	})

	t.Run("NotExistsID", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		stmt := `
		INSERT INTO
		  passwords (name, data, created_at, updated_at)
		VALUES (?, ?, ?, ?);
		`
		name := "testName"
		data := []byte("helloTestWorld")
		createdAt := time.Now()
		updatedAt := createdAt
		_, err := st.s.ExecContext(
			st.ctx, stmt, name, data, createdAt, updatedAt,
		)
		require.NoError(t, err)

		id := 2
		retData, err := st.r.ReadByID(st.ctx, id)
		require.ErrorIs(t, err, repository.ErrNotExists)
		assert.Nil(t, retData)
	})

	t.Run("ExistsIDButDeleted", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		stmt := `
		INSERT INTO
		  passwords (name, data, created_at, updated_at, deleted)
		VALUES (?, ?, ?, ?, ?);
		`
		name := "testName"
		data := []byte("helloTestWorld")
		createdAt := time.Now()
		updatedAt := createdAt
		deleted := true
		_, err := st.s.ExecContext(
			st.ctx, stmt, name, data, createdAt, updatedAt, deleted,
		)
		require.NoError(t, err)

		id := 1
		retData, err := st.r.ReadByID(st.ctx, id)
		require.ErrorIs(t, err, repository.ErrNotExists)
		assert.Nil(t, retData)
	})
}

func TestListNames(t *testing.T) {
	type obj struct {
		name      string
		data      []byte
		createdAt time.Time
		updatedAt time.Time
		deleted   bool
	}

	t.Run("Ordinary", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		stmt, err := st.s.PrepareContext(
			st.ctx,
			`INSERT INTO
			  passwords (name, data, created_at, updated_at, deleted)
			VALUES (?, ?, ?, ?, ?);`,
		)
		require.NoError(t, err)

		tNow := time.Now()
		notDeleted := false
		inserts := []obj{
			{"A", []byte("A"), tNow, tNow, notDeleted},
			{"C", []byte("C"), tNow, tNow, notDeleted},
			{"B", []byte("B"), tNow, tNow, notDeleted},
		}

		for _, obj := range inserts {
			stmt.ExecContext(
				st.ctx,
				obj.name, obj.data, obj.createdAt, obj.updatedAt, obj.deleted,
			)
		}
		err = stmt.Close()
		require.NoError(t, err)

		data, err := st.r.ListNames(st.ctx)
		require.NoError(t, err)
		require.Len(t, data, len(inserts))

		for i, item := range data {
			id, name := item[0], item[1]
			switch i {
			case 0:
				assert.Equal(t, "1", id)
				assert.Equal(t, "A", name)
			case 1:
				assert.Equal(t, "3", id)
				assert.Equal(t, "B", name)
			case 2:
				assert.Equal(t, "2", id)
				assert.Equal(t, "C", name)
			}
		}
	})

	t.Run("HaveDeleted", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		stmt, err := st.s.PrepareContext(
			st.ctx,
			`INSERT INTO
			  passwords (name, data, created_at, updated_at, deleted)
			VALUES (?, ?, ?, ?, ?);`,
		)
		require.NoError(t, err)

		tNow := time.Now()
		notDeleted := false
		deleted := true
		inserts := []obj{
			{"A", []byte("A"), tNow, tNow, notDeleted},
			{"C", []byte("C"), tNow, tNow, deleted},
			{"B", []byte("B"), tNow, tNow, notDeleted},
		}

		for _, obj := range inserts {
			stmt.ExecContext(
				st.ctx,
				obj.name, obj.data, obj.createdAt, obj.updatedAt, obj.deleted,
			)
		}
		err = stmt.Close()
		require.NoError(t, err)

		data, err := st.r.ListNames(st.ctx)
		require.NoError(t, err)
		require.Len(t, data, len(inserts)-1)

		for i, item := range data {
			id, name := item[0], item[1]
			switch i {
			case 0:
				assert.Equal(t, "1", id)
				assert.Equal(t, "A", name)
			case 1:
				assert.Equal(t, "3", id)
				assert.Equal(t, "B", name)
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)

		insertName := "insertName"
		insertData := []byte("insertData")
		insertTime := time.Now()
		_, err := st.s.ExecContext(st.ctx,
			`
			INSERT INTO passwords (name, data, created_at, updated_at)
			VALUES
			  (?, ?, ?, ?);
			`,
			insertName, insertData, insertTime, insertTime,
		)
		require.NoError(t, err)

		entryNum := 1
		updateName := "updateName"
		updateData := []byte("updateData")
		err = st.r.Update(st.ctx, entryNum, updateName, updateData)
		require.NoError(t, err)

		var actualName string
		var actualData []byte
		var actualUpdatedAt time.Time

		err = st.s.QueryRowContext(
			st.ctx, `SELECT name, data, updated_at FROM passwords WHERE id=1;`,
		).Scan(&actualName, &actualData, &actualUpdatedAt)
		require.NoError(t, err)
		assert.Equal(t, updateName, actualName)
		assert.Equal(t, updateData, actualData)
		assert.Less(t, insertTime, actualUpdatedAt)
	})

	t.Run("NotExistsEntry", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		entryNum := 1
		updateName := "updateName"
		updateData := []byte("updateData")
		err := st.r.Update(st.ctx, entryNum, updateName, updateData)
		assert.ErrorIs(t, err, repository.ErrNotExists)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Ordinary", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)

		insertName := "insertName"
		insertData := []byte("insertData")
		insertTime := time.Now()
		_, err := st.s.ExecContext(st.ctx,
			`
			INSERT INTO passwords (name, data, created_at, updated_at)
			VALUES
			  (?, ?, ?, ?);
			`,
			insertName, insertData, insertTime, insertTime,
		)
		require.NoError(t, err)

		entryNum := 1
		err = st.r.Delete(st.ctx, entryNum)
		require.NoError(t, err)

		var actualName sql.NullString
		var actualData []byte
		var actualUpdatedAt time.Time
		var actualDeleted bool
		err = st.s.QueryRowContext(
			st.ctx, `SELECT name, data, updated_at, deleted FROM passwords WHERE id=1;`,
		).Scan(&actualName, &actualData, &actualUpdatedAt, &actualDeleted)
		require.NoError(t, err)

		var expectedData []byte
		assert.False(t, actualName.Valid)
		assert.Equal(t, expectedData, actualData)
		assert.True(t, actualDeleted)
		assert.Less(t, insertTime, actualUpdatedAt)
	})

	t.Run("NotExistsEntry", func(t *testing.T) {
		st := newSuite(t, repository.NewPwd)
		entryNum := 1
		err := st.r.Delete(st.ctx, entryNum)
		assert.ErrorIs(t, err, repository.ErrNotExists)
	})
}

func newEntitySyncSuite(t *testing.T, newRepoFn repoConstructor) *suite {
	log := logger.NewPretty("debug")
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx := t.Context()

	t.Cleanup(func() {
		os.Remove(dsn)
	})

	s := storage.New(log, dsn)
	s.MustRun(ctx)
	r := newRepoFn(log, s)

	return &suite{ctx, r, s}
}
