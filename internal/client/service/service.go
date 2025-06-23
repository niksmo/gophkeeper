package service

import "errors"

var (
	ErrAlreadyExists = errors.New("object already exists")
	ErrNotExists     = errors.New("object not exists")
	ErrInvalidKey    = errors.New("invalid key provided")
)
