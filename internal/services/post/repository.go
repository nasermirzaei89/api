package post

import (
	"context"
)

type Repository interface {
	FindBySlug(ctx context.Context, slug string) (res *Entity, err error)
	Insert(ctx context.Context, entity Entity) (err error)
}
