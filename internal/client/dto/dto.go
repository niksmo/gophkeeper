package dto

import "time"

// * Storage *

type (
	PWD struct {
		Name     string
		Login    string
		Password string
	}

	BIN struct {
		Name string
		Ext  string
		Data []byte
	}

	BankCard struct {
		Name       string
		Number     string
		ExpDate    string
		HolderName string
	}

	Text struct {
		Name string
		Data string
	}

	Sync struct {
		ID        int
		PID       int
		StartedAt time.Time
		StoppedAt *time.Time
	}
)

// * Synchronization *

type (
	SyncComparable struct {
		ID        int
		Name      string
		UpdatedAt time.Time
	}

	SyncPayload struct {
		ID        int
		Name      string
		Data      []byte
		CreatedAt time.Time
		UpdatedAt time.Time
		Deleted   bool
	}

	ServerComparable = SyncComparable
	ServerPayload    = SyncPayload

	LocalComparable struct {
		SyncComparable
		SyncID int
	}

	LocalPayload struct {
		SyncPayload
		SyncID int
	}
)
