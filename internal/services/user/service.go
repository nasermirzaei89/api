package user

import "context"

type Service interface {
	LogIn(ctx context.Context, req LogInRequest) (res *LogInResponse, err error)
	GetUserByUUID(ctx context.Context, userID string) (res *Entity, err error)
	GetUserByTokenString(ctx context.Context, tokenString string) (res *Entity, err error)
}

type LogInRequest struct {
	Username string
	Password string
}

type LogInResponse struct {
	AccessToken string
	UserUUID    string
}
