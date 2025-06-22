package dto

import "time"

type User struct {
	ID        int
	Login     string
	Password  string
	CreatedAt time.Time
	Disabled  bool
}
