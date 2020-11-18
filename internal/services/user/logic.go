package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/nasermirzaei89/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	repo            Repository
	signKey         []byte
	verificationKey []byte
}

func (svc *service) GetUserByTokenString(ctx context.Context, tokenString string) (*Entity, error) {
	err := jwt.Verify(tokenString, svc.verificationKey)
	if err != nil {
		return nil, errors.Wrap(err, "error on verify jwt token")
	}

	token, err := jwt.Parse(tokenString)
	if err != nil {
		return nil, errors.Wrap(err, "error on parse jwt token")
	}

	subject, err := token.GetSubject()
	if err != nil {
		return nil, errors.Wrap(err, "error on get token subject")
	}

	entity, err := svc.repo.FindByUUID(ctx, subject)
	if err != nil {
		return nil, errors.Wrap(err, "error on find user by uuid")
	}

	if entity == nil {
		return nil, ErrUserWithUUIDNotFound{UUID: subject}
	}

	return entity, nil
}

func (svc *service) GetUserByUUID(ctx context.Context, userUUID string) (*Entity, error) {
	entity, err := svc.repo.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, errors.Wrap(err, "error on find by uuid")
	}

	if entity == nil {
		return nil, ErrUserWithUUIDNotFound{UUID: userUUID}
	}

	return entity, nil
}

func (svc *service) LogIn(ctx context.Context, req LogInRequest) (*LogInResponse, error) {
	entity, err := svc.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "error on find by username")
	}

	if entity == nil {
		return nil, ErrUserWithUsernameNotFound{Username: req.Username}
	}

	err = bcrypt.CompareHashAndPassword([]byte(entity.PasswordHash), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidPasswordReceived{}
		}

		return nil, errors.Wrap(err, "error on compare hash and password")
	}

	token := jwt.New(jwt.RS256)
	token.SetSubject(entity.UUID)
	token.SetIssuedAt(time.Now())
	token.SetJWTID(uuid.New().String())

	accessToken, err := jwt.Sign(token, svc.signKey)
	if err != nil {
		return nil, errors.Wrap(err, "error on sign token")
	}

	rsp := LogInResponse{
		AccessToken: accessToken,
		UserUUID:    entity.UUID,
	}

	return &rsp, nil
}

func NewService(repo Repository, signKey, verificationKey []byte) Service {
	svc := service{
		repo:            repo,
		signKey:         signKey,
		verificationKey: verificationKey,
	}

	return &svc
}
