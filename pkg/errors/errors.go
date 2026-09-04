package errors

import "errors"

var (
	ErrAuthNotValid                = errors.New("auth not valid")
	ErrInvalidAuthorizationHeader  = errors.New("invalid authorization header")
	ErrTokenIsExpired              = errors.New("token is expired")
	ErrTokenisBelongsToAnotherUser = errors.New("token is belongs to another user")
	ErrLoginBusy                   = errors.New("login is busy")
	ErrNotFound                    = errors.New("not found")
	ErrWriteToSupport              = errors.New("error, please write to support")
	ErrTokenIsNotValid             = errors.New("token is not valid")
	ErrInvalidOrderNumberFormat    = errors.New("invalid order number format")
	ErrEmptyListOrders             = errors.New("empty list orders")
	ErrLessDrawals                 = errors.New("drawals less then sum")
)
