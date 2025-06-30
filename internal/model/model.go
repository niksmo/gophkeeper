package model

import (
	"database/sql"
	"time"
)

type Row interface {
	Scan(dest ...any) error
}

type SyncComparable struct {
	ID        int64
	Name      string
	UpdatedAt time.Time
}

func (sc *SyncComparable) ScanRow(row Row) error {
	return row.Scan(&sc.ID, &sc.Name, &sc.UpdatedAt)
}

type SyncPayload struct {
	ID        int64
	Name      string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
}

func (sp *SyncPayload) ScanRow(row Row) error {
	return row.Scan(&sp.ID, &sp.Name, &sp.Data,
		&sp.CreatedAt, &sp.UpdatedAt, &sp.Deleted)
}

type LocalComparable struct {
	SyncComparable
	SyncID int64
}

func (lc *LocalComparable) ScanRow(row Row) error {
	var syncID sql.NullInt64
	if err := row.Scan(&lc.ID, &lc.Name, &lc.UpdatedAt, &syncID); err != nil {
		return err
	}
	lc.SyncID = syncID.Int64
	return nil
}

type LocalPayload struct {
	SyncPayload
	SyncID int64
}

func (lp *LocalPayload) ScanRow(row Row) error {
	var syncID sql.NullInt64
	err := row.Scan(&lp.ID, &lp.Name, &lp.Data, &lp.CreatedAt,
		&lp.UpdatedAt, &lp.Deleted, &syncID)
	if err != nil {
		return err
	}
	lp.SyncID = syncID.Int64
	return nil
}
