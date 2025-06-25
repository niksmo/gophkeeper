package synchronization

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

const (
	GetLtData    = "GetLtData"
	GetAllData   = "GetAllData"
	GetData      = "GetData"
	InsertData   = "InsertData"
	InsertSyncID = "InsertSyncID"
)

type MockServer struct {
	mock.Mock
}

func (s *MockServer) GetLtData() []LtSrvDTO {
	args := s.Called()
	return args.Get(0).([]LtSrvDTO)
}

func (s *MockServer) GetAllData() []SrvDTO {
	args := s.Called()
	return args.Get(0).([]SrvDTO)
}

func (s *MockServer) GetData(idS []int) []SrvDTO {
	args := s.Called(idS)
	return args.Get(0).([]SrvDTO)
}

func (s *MockServer) InsertData(d []ClDTO) []int {
	args := s.Called(d)
	return args.Get(0).([]int)
}

type MockClient struct {
	mock.Mock
}

func (c *MockClient) GetLtData() []LtClDTO {
	args := c.Called()
	return args.Get(0).([]LtClDTO)
}

func (c *MockClient) GetAllData() []ClDTO {
	args := c.Called()
	return args.Get(0).([]ClDTO)
}

func (c *MockClient) GetData(idS []int) []ClDTO {
	args := c.Called(idS)
	return args.Get(0).([]ClDTO)
}

func (c *MockClient) InsertData(d []SrvDTO) {
	c.Called(d)
}

func (c *MockClient) InsertSyncID(kv [][2]int) {
	c.Called(kv)
}

func TestNoServerData(t *testing.T) {
	defer prettyPanic(t)
	srv := &MockServer{}
	cl := &MockClient{}

	clData := []ClDTO{
		{
			ID:        1,
			Name:      "A",
			Data:      []byte("dataA"),
			CreatedAt: getDate("20/07/2025 18:00"),
			UpdatedAt: getDate("20/07/2025 18:00"),
			Deleted:   false,
		},
		{
			ID:        2,
			Name:      "B",
			Data:      []byte("dataB"),
			CreatedAt: getDate("20/07/2025 19:00"),
			UpdatedAt: getDate("20/07/2025 19:00"),
			Deleted:   false,
		},
	}

	srv.On(GetLtData).Return([]LtSrvDTO{})

	cl.On(GetAllData).Return(clData)

	srv.On(InsertData, clData).Return([]int{1, 2})

	cl.On(InsertSyncID, [][2]int{{1, 1}, {2, 2}})

	SyncWithUniqueName(srv, cl)

	srv.AssertNumberOfCalls(t, GetLtData, 1)
	cl.AssertNumberOfCalls(t, GetAllData, 1)
	srv.AssertNumberOfCalls(t, InsertData, 1)
	cl.AssertNumberOfCalls(t, InsertSyncID, 1)
}

func TestNoServerAndClientData(t *testing.T) {
	defer prettyPanic(t)
	srv := &MockServer{}
	cl := &MockClient{}

	srv.On(GetLtData).Return([]LtSrvDTO{})

	cl.On(GetAllData).Return([]ClDTO{})

	SyncWithUniqueName(srv, cl)

	srv.AssertNumberOfCalls(t, GetLtData, 1)
	cl.AssertNumberOfCalls(t, GetAllData, 1)
	srv.AssertNumberOfCalls(t, InsertData, 0)
	cl.AssertNumberOfCalls(t, InsertSyncID, 0)
}

func TestNoClientData(t *testing.T) {
	defer prettyPanic(t)
	srv := &MockServer{}
	cl := &MockClient{}

	srvLtData := []LtSrvDTO{
		{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 18:00")},
		{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 19:00")},
	}

	srvData := []SrvDTO{
		{
			ID:        1,
			Name:      "A",
			Data:      []byte("dataA"),
			CreatedAt: getDate("20/07/2025 18:00"),
			UpdatedAt: getDate("20/07/2025 18:00"),
			Deleted:   false,
		},
		{
			ID:        2,
			Name:      "B",
			Data:      []byte("dataB"),
			CreatedAt: getDate("20/07/2025 19:00"),
			UpdatedAt: getDate("20/07/2025 19:00"),
			Deleted:   false,
		},
	}

	clLtData := []LtClDTO{}

	srv.On(GetLtData).Return(srvLtData)

	cl.On(GetLtData).Return(clLtData)

	srv.On(GetAllData).Return(srvData)

	cl.On(InsertData, srvData)

	SyncWithUniqueName(srv, cl)

	srv.AssertNumberOfCalls(t, GetLtData, 1)
	cl.AssertNumberOfCalls(t, GetLtData, 1)
	srv.AssertNumberOfCalls(t, GetAllData, 1)
	cl.AssertNumberOfCalls(t, InsertData, 1)
}

func TestDifferentData(t *testing.T) {
	t.Run("EqualLenNoSyncID", func(t *testing.T) {
		defer prettyPanic(t)
		srv := &MockServer{}
		cl := &MockClient{}

		srvLtData := []LtSrvDTO{
			{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 18:00")},
			{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 19:00")},
		}

		srvData := []SrvDTO{
			{
				ID:        1,
				Name:      "A",
				Data:      []byte("dataA"),
				CreatedAt: getDate("20/07/2025 18:00"),
				UpdatedAt: getDate("20/07/2025 18:00"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "B",
				Data:      []byte("dataB"),
				CreatedAt: getDate("20/07/2025 19:00"),
				UpdatedAt: getDate("20/07/2025 19:00"),
				Deleted:   false,
			},
		}

		clLtData := []LtClDTO{
			{ID: 1, Name: "C", UpdatedAt: getDate("20/07/2025 18:01"), SyncID: 0},
			{ID: 2, Name: "D", UpdatedAt: getDate("20/07/2025 19:01"), SyncID: 0},
		}

		clData := []ClDTO{
			{
				ID:        1,
				Name:      "C",
				Data:      []byte("dataC"),
				CreatedAt: getDate("20/07/2025 18:01"),
				UpdatedAt: getDate("20/07/2025 18:01"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "D",
				Data:      []byte("dataD"),
				CreatedAt: getDate("20/07/2025 19:01"),
				UpdatedAt: getDate("20/07/2025 19:01"),
				Deleted:   false,
			},
		}

		srv.On(GetLtData).Return(srvLtData)
		cl.On(GetLtData).Return(clLtData)
		srv.On(GetData, []int{1, 2}).Return(srvData)
		cl.On(InsertData, srvData)
		cl.On(GetData, []int{1, 2}).Return(clData)
		srv.On(InsertData, clData).Return([]int{3, 4})
		cl.On(InsertSyncID, [][2]int{{1, 3}, {2, 4}})

		SyncWithUniqueName(srv, cl)

		srv.AssertNumberOfCalls(t, GetLtData, 1)
		srv.AssertNumberOfCalls(t, GetData, 1)
		srv.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, GetLtData, 1)
		cl.AssertNumberOfCalls(t, GetData, 1)
		cl.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, InsertSyncID, 1)
	})

	t.Run("ClDataMoreThenSrvNoSyncID", func(t *testing.T) {
		defer prettyPanic(t)
		srv := &MockServer{}
		cl := &MockClient{}

		srvLtData := []LtSrvDTO{
			{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 18:00")},
			{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 19:00")},
		}

		srvData := []SrvDTO{
			{
				ID:        1,
				Name:      "A",
				Data:      []byte("dataA"),
				CreatedAt: getDate("20/07/2025 18:00"),
				UpdatedAt: getDate("20/07/2025 18:00"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "B",
				Data:      []byte("dataB"),
				CreatedAt: getDate("20/07/2025 19:00"),
				UpdatedAt: getDate("20/07/2025 19:00"),
				Deleted:   false,
			},
		}

		clLtData := []LtClDTO{
			{ID: 1, Name: "C", UpdatedAt: getDate("20/07/2025 18:01"), SyncID: 0},
			{ID: 2, Name: "D", UpdatedAt: getDate("20/07/2025 19:01"), SyncID: 0},
			{ID: 3, Name: "E", UpdatedAt: getDate("20/07/2025 19:01"), SyncID: 0},
		}

		clData := []ClDTO{
			{
				ID:        1,
				Name:      "C",
				Data:      []byte("dataC"),
				CreatedAt: getDate("20/07/2025 18:01"),
				UpdatedAt: getDate("20/07/2025 18:01"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "D",
				Data:      []byte("dataD"),
				CreatedAt: getDate("20/07/2025 19:01"),
				UpdatedAt: getDate("20/07/2025 19:01"),
				Deleted:   false,
			},
			{
				ID:        3,
				Name:      "E",
				Data:      []byte("dataE"),
				CreatedAt: getDate("20/07/2025 20:01"),
				UpdatedAt: getDate("20/07/2025 20:01"),
				Deleted:   false,
			},
		}

		srv.On(GetLtData).Return(srvLtData)
		cl.On(GetLtData).Return(clLtData)
		srv.On(GetData, []int{1, 2}).Return(srvData)
		cl.On(InsertData, srvData)
		cl.On(GetData, []int{1, 2, 3}).Return(clData)
		srv.On(InsertData, clData).Return([]int{3, 4, 5})
		cl.On(InsertSyncID, [][2]int{{1, 3}, {2, 4}, {3, 5}})

		SyncWithUniqueName(srv, cl)

		srv.AssertNumberOfCalls(t, GetLtData, 1)
		srv.AssertNumberOfCalls(t, GetData, 1)
		srv.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, GetLtData, 1)
		cl.AssertNumberOfCalls(t, GetData, 1)
		cl.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, InsertSyncID, 1)
	})

	t.Run("SrvDataMoreThenClNoSyncID", func(t *testing.T) {
		defer prettyPanic(t)
		srv := &MockServer{}
		cl := &MockClient{}

		srvLtData := []LtSrvDTO{
			{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 18:00")},
			{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 19:00")},
			{ID: 3, Name: "E", UpdatedAt: getDate("20/07/2025 20:00")},
		}

		srvData := []SrvDTO{
			{
				ID:        1,
				Name:      "A",
				Data:      []byte("dataA"),
				CreatedAt: getDate("20/07/2025 18:00"),
				UpdatedAt: getDate("20/07/2025 18:00"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "B",
				Data:      []byte("dataB"),
				CreatedAt: getDate("20/07/2025 19:00"),
				UpdatedAt: getDate("20/07/2025 19:00"),
				Deleted:   false,
			},
			{
				ID:        3,
				Name:      "E",
				Data:      []byte("dataE"),
				CreatedAt: getDate("20/07/2025 20:00"),
				UpdatedAt: getDate("20/07/2025 20:00"),
				Deleted:   false,
			},
		}

		clLtData := []LtClDTO{
			{ID: 1, Name: "C", UpdatedAt: getDate("20/07/2025 18:01"), SyncID: 0},
			{ID: 2, Name: "D", UpdatedAt: getDate("20/07/2025 19:01"), SyncID: 0},
		}

		clData := []ClDTO{
			{
				ID:        1,
				Name:      "C",
				Data:      []byte("dataC"),
				CreatedAt: getDate("20/07/2025 18:01"),
				UpdatedAt: getDate("20/07/2025 18:01"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "D",
				Data:      []byte("dataD"),
				CreatedAt: getDate("20/07/2025 19:01"),
				UpdatedAt: getDate("20/07/2025 19:01"),
				Deleted:   false,
			},
		}

		srv.On(GetLtData).Return(srvLtData)
		cl.On(GetLtData).Return(clLtData)
		srv.On(GetData, []int{1, 2, 3}).Return(srvData)
		cl.On(InsertData, srvData)
		cl.On(GetData, []int{1, 2}).Return(clData)
		srv.On(InsertData, clData).Return([]int{4, 5})
		cl.On(InsertSyncID, [][2]int{{1, 4}, {2, 5}})

		SyncWithUniqueName(srv, cl)

		srv.AssertNumberOfCalls(t, GetLtData, 1)
		srv.AssertNumberOfCalls(t, GetData, 1)
		srv.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, GetLtData, 1)
		cl.AssertNumberOfCalls(t, GetData, 1)
		cl.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, InsertSyncID, 1)
	})
}

func TestEqualNamesData(t *testing.T) {
	t.Run("SrvUpdateAfterClNoSyncID", func(t *testing.T) {
		defer prettyPanic(t)
		srv := &MockServer{}
		cl := &MockClient{}

		srvLtData := []LtSrvDTO{
			{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 18:00")},
			{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 20:00")},
		}

		srvData := []SrvDTO{
			{
				ID:        1,
				Name:      "A",
				Data:      []byte("dataA"),
				CreatedAt: getDate("20/07/2025 18:00"),
				UpdatedAt: getDate("20/07/2025 18:00"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "B",
				Data:      []byte("dataB"),
				CreatedAt: getDate("20/07/2025 19:00"),
				UpdatedAt: getDate("20/07/2025 20:00"),
				Deleted:   false,
			},
		}

		clLtData := []LtClDTO{
			{ID: 1, Name: "B", UpdatedAt: getDate("20/07/2025 18:01"), SyncID: 0},
			{ID: 2, Name: "C", UpdatedAt: getDate("20/07/2025 19:01"), SyncID: 0},
		}

		clData := []ClDTO{
			{
				ID:        2,
				Name:      "C",
				Data:      []byte("dataC"),
				CreatedAt: getDate("20/07/2025 19:01"),
				UpdatedAt: getDate("20/07/2025 19:01"),
				Deleted:   false,
			},
		}

		srv.On(GetLtData).Return(srvLtData)
		cl.On(GetLtData).Return(clLtData)
		srv.On(GetData, []int{1, 2}).Return(srvData)
		cl.On(InsertData, srvData)
		cl.On(GetData, []int{2}).Return(clData)
		srv.On(InsertData, clData).Return([]int{3})
		cl.On(InsertSyncID, [][2]int{{2, 3}})

		SyncWithUniqueName(srv, cl)

		srv.AssertNumberOfCalls(t, GetLtData, 1)
		srv.AssertNumberOfCalls(t, GetData, 1)
		srv.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, GetLtData, 1)
		cl.AssertNumberOfCalls(t, GetData, 1)
		cl.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, InsertSyncID, 1)
	})

	t.Run("SrvUpdateBeforeClNoSyncID", func(t *testing.T) {
		defer prettyPanic(t)
		srv := &MockServer{}
		cl := &MockClient{}

		srvLtData := []LtSrvDTO{
			{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 18:00")},
			{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 19:00")},
		}

		srvData := []SrvDTO{
			{
				ID:        1,
				Name:      "A",
				Data:      []byte("dataA"),
				CreatedAt: getDate("20/07/2025 18:00"),
				UpdatedAt: getDate("20/07/2025 18:00"),
				Deleted:   false,
			},
		}

		clLtData := []LtClDTO{
			{ID: 1, Name: "B", UpdatedAt: getDate("20/07/2025 19:01"), SyncID: 0},
			{ID: 2, Name: "C", UpdatedAt: getDate("20/07/2025 20:01"), SyncID: 0},
		}

		clData := []ClDTO{
			{
				ID:        1,
				Name:      "B",
				Data:      []byte("dataB"),
				CreatedAt: getDate("20/07/2025 19:01"),
				UpdatedAt: getDate("20/07/2025 19:01"),
				Deleted:   false,
			},
			{
				ID:        2,
				Name:      "C",
				Data:      []byte("dataC"),
				CreatedAt: getDate("20/07/2025 20:01"),
				UpdatedAt: getDate("20/07/2025 20:01"),
				Deleted:   false,
			},
		}

		srv.On(GetLtData).Return(srvLtData)
		cl.On(GetLtData).Return(clLtData)
		srv.On(GetData, []int{1}).Return(srvData)
		cl.On(InsertData, srvData)
		cl.On(GetData, []int{1, 2}).Return(clData)
		srv.On(InsertData, clData).Return([]int{2, 3})
		cl.On(InsertSyncID, [][2]int{{1, 2}, {2, 3}})

		SyncWithUniqueName(srv, cl)

		srv.AssertNumberOfCalls(t, GetLtData, 1)
		srv.AssertNumberOfCalls(t, GetData, 1)
		srv.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, GetLtData, 1)
		cl.AssertNumberOfCalls(t, GetData, 1)
		cl.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, InsertSyncID, 1)
	})

	t.Run("SrvUpdateEqualClHaveSyncedData", func(t *testing.T) {
		defer prettyPanic(t)
		srv := &MockServer{}
		cl := &MockClient{}

		srvLtData := []LtSrvDTO{
			{ID: 1, Name: "C", UpdatedAt: getDate("20/07/2025 19:00")},
			{ID: 2, Name: "A", UpdatedAt: getDate("20/07/2025 17:00")},
			{ID: 3, Name: "B", UpdatedAt: getDate("20/07/2025 18:00")},
		}

		clLtData := []LtClDTO{
			{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 17:00"), SyncID: 2},
			{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 18:00"), SyncID: 3},
			{ID: 3, Name: "C", UpdatedAt: getDate("20/07/2025 19:00"), SyncID: 1},
			{ID: 4, Name: "D", UpdatedAt: getDate("20/07/2025 20:00"), SyncID: 0},
		}

		clData := []ClDTO{
			{
				ID:        4,
				Name:      "D",
				Data:      []byte("dataD"),
				CreatedAt: getDate("20/07/2025 20:00"),
				UpdatedAt: getDate("20/07/2025 20:00"),
				Deleted:   false,
				SyncID:    0,
			},
		}

		srv.On(GetLtData).Return(srvLtData)
		cl.On(GetLtData).Return(clLtData)
		cl.On(GetData, []int{4}).Return(clData)
		srv.On(InsertData, clData).Return([]int{4})
		cl.On(InsertSyncID, [][2]int{{4, 4}})

		SyncWithUniqueName(srv, cl)

		srv.AssertNumberOfCalls(t, GetLtData, 1)
		srv.AssertNumberOfCalls(t, GetData, 0)
		srv.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, GetLtData, 1)
		cl.AssertNumberOfCalls(t, GetData, 1)
		cl.AssertNumberOfCalls(t, InsertData, 0)
		cl.AssertNumberOfCalls(t, InsertSyncID, 1)
	})

	t.Run("SrvUpdateNotEqualClHaveSyncedData", func(t *testing.T) {
		defer prettyPanic(t)
		srv := &MockServer{}
		cl := &MockClient{}

		srvLtData := []LtSrvDTO{
			{ID: 1, Name: "C", UpdatedAt: getDate("20/07/2025 19:00")},
			{ID: 2, Name: "A", UpdatedAt: getDate("20/07/2025 18:00")},
			{ID: 3, Name: "B", UpdatedAt: getDate("20/07/2025 18:00")},
			{ID: 4, Name: "E", UpdatedAt: getDate("20/07/2025 21:00")},
		}

		clLtData := []LtClDTO{
			{ID: 1, Name: "A", UpdatedAt: getDate("20/07/2025 17:00"), SyncID: 2},
			{ID: 2, Name: "B", UpdatedAt: getDate("20/07/2025 18:00"), SyncID: 3},
			{ID: 3, Name: "C", UpdatedAt: getDate("20/07/2025 20:00"), SyncID: 1},
			{ID: 4, Name: "D", UpdatedAt: getDate("20/07/2025 22:00"), SyncID: 0},
		}

		srvData := []SrvDTO{
			{
				ID:        2,
				Name:      "A",
				Data:      []byte("dataA"),
				CreatedAt: getDate("20/07/2025 18:00"),
				UpdatedAt: getDate("20/07/2025 18:00"),
			},
			{
				ID:        4,
				Name:      "E",
				Data:      []byte("dataE"),
				CreatedAt: getDate("20/07/2025 21:00"),
				UpdatedAt: getDate("20/07/2025 21:00"),
			},
		}

		clData := []ClDTO{
			{
				ID:        3,
				Name:      "C",
				Data:      []byte("dataC"),
				CreatedAt: getDate("20/07/2025 20:00"),
				UpdatedAt: getDate("20/07/2025 20:00"),
				SyncID:    1,
			},
			{
				ID:        4,
				Name:      "D",
				Data:      []byte("dataD"),
				CreatedAt: getDate("20/07/2025 22:00"),
				UpdatedAt: getDate("20/07/2025 22:00"),
				SyncID:    0,
			},
		}

		srv.On(GetLtData).Return(srvLtData)
		cl.On(GetLtData).Return(clLtData)
		srv.On(GetData, []int{2, 4}).Return(srvData)
		cl.On(InsertData, srvData)
		cl.On(GetData, []int{3, 4}).Return(clData)
		srv.On(InsertData, clData).Return([]int{1, 5})
		cl.On(InsertSyncID, [][2]int{{4, 5}})

		SyncWithUniqueName(srv, cl)

		srv.AssertNumberOfCalls(t, GetLtData, 1)
		srv.AssertNumberOfCalls(t, GetData, 1)
		srv.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, GetLtData, 1)
		cl.AssertNumberOfCalls(t, GetData, 1)
		cl.AssertNumberOfCalls(t, InsertData, 1)
		cl.AssertNumberOfCalls(t, InsertSyncID, 1)
	})
}

// Layout is "02/01/2006 15:04"
func getDate(d string) time.Time {
	const layout = "02/01/2006 15:04"
	date, err := time.Parse(layout, d)
	if err != nil {
		panic(err)
	}
	return date
}

func prettyPanic(t *testing.T) {
	t.Helper()
	if r := recover(); r != nil {
		t.Log(r)
		t.Fail()
	}
}
