package post

import (
	"context"
)

type Repository interface {
	FindBySlug(ctx context.Context, slug string) (res *Entity, err error)
	Insert(ctx context.Context, entity Entity) (err error)
	FindByUUID(ctx context.Context, uuid string) (res *Entity, err error)
	List(ctx context.Context) (res []*Entity, err error)
}
