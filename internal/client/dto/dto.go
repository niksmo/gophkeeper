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
	Name       string
	Number     string
	ExpDate    string
	HolderName string
}

type Text struct {
	Name string
	Data string
}

type Sync struct {
	ID        int
	PID       int
	StartedAt time.Time
	StoppedAt *time.Time
}

// * Synchronization DTO's *

// Embedded
type (
	comparable struct {
		ID        int
		Name      string
		UpdatedAt time.Time
	}

	payload struct {
		ID        int
		Name      string
		Data      []byte
		CreatedAt time.Time
		UpdatedAt time.Time
		Deleted   bool
	}
)

type ServerComparable struct {
	comparable
}

type LocalComparable struct {
	comparable
	SyncID int
}

type ServerPayload struct {
	payload
}

type LocalPayload struct {
	payload
	SyncID int
}
