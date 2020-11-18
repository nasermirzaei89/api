package post

import (
	"context"
)

type Service interface {
	CreatePost(ctx context.Context, req CreatePostRequest) (res *Entity, err error)
	GetPostByUUID(ctx context.Context, postUUID string) (res *Entity, err error)
}

type CreatePostRequest struct {
	Title           string
	Slug            string
	ContentMarkdown string
}
