package domain

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrBadRequest          = errors.New("bad request")
	ErrInvalidParam        = errors.New("invalid param")
	ErrInternalServerError = errors.New("internal server error")
)
