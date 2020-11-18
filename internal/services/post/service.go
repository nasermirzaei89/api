package post

import (
	"context"
)

type Service interface {
	CreatePost(ctx context.Context, req CreatePostRequest) (res *Entity, err error)
	GetPostByUUID(ctx context.Context, postUUID string) (res *Entity, err error)
	GetPostBySlug(ctx context.Context, slug string) (res *Entity, err error)
	ListPosts(ctx context.Context) (res []*Entity, err error)
	UpdatePostByUUID(ctx context.Context, postUUID string, req UpdatePostByUUIDRequest) (res *Entity, err error)
}

type CreatePostRequest struct {
	Title           string
	Slug            string
	ContentMarkdown string
}

type UpdatePostByUUIDRequest struct {
	Title           string
	Slug            string
	ContentMarkdown string
}
