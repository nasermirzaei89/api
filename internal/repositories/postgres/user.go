package postgres

import (
	"context"
	"database/sql"
	"github.com/nasermirzaei89/api/internal/services/user"
	"github.com/pkg/errors"
)

type userRepo struct {
	db *sql.DB
}

func (repo *userRepo) FindByUsername(ctx context.Context, username string) (*user.Entity, error) {
	var model user.Entity

	// prepare query
	query := `SELECT uuid, username, password_hash FROM users WHERE username = $1;`
	args := []interface{}{username}
	dest := []interface{}{&model.UUID, &model.Username, &model.PasswordHash}

	if err := repo.db.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(errors.WithStack(err), "error on query row")
	}

	return &model, nil
}

func (repo *userRepo) FindByUUID(ctx context.Context, userUUID string) (*user.Entity, error) {
	var model user.Entity

	// prepare query
	query := `SELECT uuid, username, password_hash FROM users WHERE uuid = $1;`
	args := []interface{}{userUUID}
	dest := []interface{}{&model.UUID, &model.Username, &model.PasswordHash}

	if err := repo.db.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(errors.WithStack(err), "error on query row")
	}

	return &model, nil
}

func NewUserRepository(db *sql.DB) user.Repository {
	repo := userRepo{
		db: db,
	}

	return &repo
}
