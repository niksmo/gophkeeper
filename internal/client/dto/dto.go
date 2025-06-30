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
