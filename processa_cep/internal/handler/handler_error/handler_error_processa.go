package handler_error

import "errors"

var (
    ErrInvalidZipcode  = errors.New("invalid zipcode")
    ErrZipcodeNotFound = errors.New("can not find zipcode")
    ErrInternal        = errors.New("internal error")
)
