package dto

import "time"

type User struct {
	ID           int
	Login        string
	PasswordHash []byte
	CreatedAt    time.Time
	Disabled     bool
}
