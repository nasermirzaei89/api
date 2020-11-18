package postgres

import (
	"context"
	"database/sql"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/pkg/errors"
)

type postRepo struct {
	db *sql.DB
}

func (repo *postRepo) UpdateByUUID(ctx context.Context, uuid string, entity post.Entity) error {
	query := `UPDATE posts SET uuid = $1, title = $2, slug = $3, content_markdown = $4, content_html = $5 WHERE uuid = $6;`
	args := []interface{}{entity.UUID, entity.Title, entity.Slug, entity.ContentMarkdown, entity.ContentHTML, uuid}

	_, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "error on exec")
	}

	return nil
}

func (repo *postRepo) List(ctx context.Context) ([]*post.Entity, error) {
	query := `SELECT uuid, title, slug, content_markdown, content_html FROM posts;`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on query")
	}

	res := make([]*post.Entity, 0)
	for rows.Next() {
		var model post.Entity
		dest := []interface{}{&model.UUID, &model.Title, &model.Slug, &model.ContentMarkdown, &model.ContentHTML}
		err = rows.Scan(dest...)
		if err != nil {
			return nil, errors.Wrap(err, "error on scan row")
		}

		res = append(res, &model)
	}

	return res, nil
}

func (repo *postRepo) FindByUUID(ctx context.Context, uuid string) (*post.Entity, error) {
	var model post.Entity

	// prepare query
	query := `SELECT uuid, title, slug, content_markdown, content_html FROM posts WHERE uuid = $1;`
	args := []interface{}{uuid}
	dest := []interface{}{&model.UUID, &model.Title, &model.Slug, &model.ContentMarkdown, &model.ContentHTML}

	err := repo.db.QueryRowContext(ctx, query, args...).Scan(dest...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(errors.WithStack(err), "error on query row")
	}

	return &model, nil
}

func (repo *postRepo) FindBySlug(ctx context.Context, slug string) (*post.Entity, error) {
	var model post.Entity

	// prepare query
	query := `SELECT uuid, title, slug, content_markdown, content_html FROM posts WHERE slug = $1;`
	args := []interface{}{slug}
	dest := []interface{}{&model.UUID, &model.Title, &model.Slug, &model.ContentMarkdown, &model.ContentHTML}

	err := repo.db.QueryRowContext(ctx, query, args...).Scan(dest...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(errors.WithStack(err), "error on query row")
	}

	return &model, nil
}

func (repo *postRepo) Insert(ctx context.Context, entity post.Entity) error {
	query := `INSERT INTO posts (uuid, title, slug, content_markdown, content_html) VALUES ($1, $2, $3, $4, $5);`
	args := []interface{}{entity.UUID, entity.Title, entity.Slug, entity.ContentMarkdown, entity.ContentHTML}

	_, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "error on exec")
	}

	return nil
}

func NewPostRepository(db *sql.DB) post.Repository {
	repo := postRepo{
		db: db,
	}

	return &repo
}
