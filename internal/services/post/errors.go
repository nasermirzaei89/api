package post

import "fmt"

type ErrPostWithUUIDNotFound struct {
	UUID string
}

func (err ErrPostWithUUIDNotFound) Error() string {
	return fmt.Sprintf("post with uuid '%s' not found", err.UUID)
}
