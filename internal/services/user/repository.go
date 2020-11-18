package user

import (
	"context"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (res *Entity, err error)
	FindByUUID(ctx context.Context, userUUID string) (res *Entity, err error)
}
