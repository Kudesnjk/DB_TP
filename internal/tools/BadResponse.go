package tools

import (
	"errors"
)

const (
	ConstInternalErrorMessage = "internal server error"
	ConstNotFoundMessage      = "not found message"
	ConstSomeMessage          = "some message"
)

var (
	ErrorUserNotFound       = errors.New("User not found")
	ErrorParentPostNotFound = errors.New("Error with parent post")
)

type BadResponse struct {
	Message string `json:"message"`
}
