package post

import "fmt"

type ErrPostWithUUIDNotFound struct {
	UUID string
}

func (err ErrPostWithUUIDNotFound) Error() string {
	return fmt.Sprintf("post with uuid '%s' not found", err.UUID)
}

type ErrPostWithSlugNotFound struct {
	Slug string
}

func (err ErrPostWithSlugNotFound) Error() string {
	return fmt.Sprintf("post with slug '%s' not found", err.Slug)
}

type ErrPostWithSlugNotPublished struct {
	Slug string
}

func (err ErrPostWithSlugNotPublished) Error() string {
	return fmt.Sprintf("post with slug '%s' not published", err.Slug)
}
