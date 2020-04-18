package shared

import (
	"errors"
	"fmt"
)

var (
	// ErrTokenFormat error
	ErrTokenFormat = errors.New("Invalid token format")
	// ErrTokenExpired error
	ErrTokenExpired = errors.New("Token is expired")

	// ErrInvalidUsernameOrPassword error
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")

	// ErrUserAlreadyExist error
	ErrUserAlreadyExist = errors.New("users already exist")
)

// ErrorUserNotFound represent user not found error
type ErrorUserNotFound struct {
	Message string
}

// Error function for ErrorUserNotFound
func (e *ErrorUserNotFound) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}
