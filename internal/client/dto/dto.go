package dto

import "time"

type PWD struct {
	Name     string
	Login    string
	Password string
}

type BIN struct {
	Name string
	Ext  string
	Data []byte
}

type BankCard struct {
	Name    string
	Number  string
	ExpDate string
	CVC     string
}

type Text struct {
	Name string
	Data string
}

type LtClientEntry struct {
	ID        int
	Name      string
	UpdatedAt time.Time
	SyncID    int
}

type ClientDTO struct {
	ID        int
	Name      string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
	SyncID    int
}

type SyncDTO struct {
	ID        int
	PID       int
	StartedAt time.Time
	StoppedAt *time.Time
}
