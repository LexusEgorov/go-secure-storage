package models

import "errors"

type User struct {
	ID       int
	Email    string
	Password string
}

type Credentials struct {
	JWT     string
	Refresh string
}

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("already registered")
var ErrUnauthorized = errors.New("token not ok")
var ErrInternal = errors.New("token not ok")
