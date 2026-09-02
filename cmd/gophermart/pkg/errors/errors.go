package errors

import "errors"

var (
	ErrLoginBusy       = errors.New("login is busy")
	ErrNotFound        = errors.New("not found")
	ErrWriteToSupport  = errors.New("error, please write to support")
	ErrTokenIsNotValid = errors.New("token is not valid")
	ErrLessDrawals     = errors.New("drawals less then sum")
)
