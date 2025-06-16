package repository

import "errors"

var (
	ErrAlreadyExist = errors.New("already exists")
	ErrNotExists    = errors.New("not exists")
)
