package user

import "fmt"

type ErrUserWithUsernameNotFound struct {
	Username string
}

func (err ErrUserWithUsernameNotFound) Error() string {
	return fmt.Sprintf("user with username '%s' not found", err.Username)
}

type ErrUserWithUUIDNotFound struct {
	UUID string
}

func (err ErrUserWithUUIDNotFound) Error() string {
	return fmt.Sprintf("user with uuid '%s' not found", err.UUID)
}

type ErrInvalidPasswordReceived struct {
}

func (err ErrInvalidPasswordReceived) Error() string {
	return fmt.Sprintf("invalid password received")
}
