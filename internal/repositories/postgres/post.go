package postgres

import (
	"context"
	"database/sql"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/pkg/errors"
	"time"
)

type postModel struct {
	UUID            string
	Title           string
	Slug            string
	ContentMarkdown string
	ContentHTML     string
	PublishedAt     sql.NullTime
}

func (m postModel) ToEntity() post.Entity {
	return post.Entity{
		UUID:            m.UUID,
		Title:           m.Title,
		Slug:            m.Slug,
		ContentMarkdown: m.ContentMarkdown,
		ContentHTML:     m.ContentHTML,
		PublishedAt: func() *time.Time {
			if m.PublishedAt.Valid {
				v := m.PublishedAt.Time
				return &v
			}

			return nil
		}(),
	}
}

func (m *postModel) FromEntity(entity post.Entity) {
	m.UUID = entity.UUID
	m.Title = entity.Title
	m.Slug = entity.Slug
	m.ContentMarkdown = entity.ContentMarkdown
	m.ContentHTML = entity.ContentHTML
	m.PublishedAt.Valid = entity.PublishedAt != nil
	if entity.PublishedAt != nil {
		m.PublishedAt.Time = *entity.PublishedAt
	}
}

type postRepo struct {
	db *sql.DB
}

func (repo *postRepo) UpdateByUUID(ctx context.Context, uuid string, entity post.Entity) error {
	m := new(postModel)
	m.FromEntity(entity)

	query := `UPDATE posts SET uuid = $1, title = $2, slug = $3, content_markdown = $4, content_html = $5, published_at = $6 WHERE uuid = $7;`
	args := []interface{}{m.UUID, m.Title, m.Slug, m.ContentMarkdown, m.ContentHTML, m.PublishedAt, uuid}

	_, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "error on exec")
	}

	return nil
}

func (repo *postRepo) List(ctx context.Context) ([]*post.Entity, error) {
	query := `SELECT uuid, title, slug, content_markdown, content_html, published_at FROM posts;`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on query")
	}

	res := make([]*post.Entity, 0)
	for rows.Next() {
		var m postModel
		dest := []interface{}{&m.UUID, &m.Title, &m.Slug, &m.ContentMarkdown, &m.ContentHTML, &m.PublishedAt}
		err = rows.Scan(dest...)
		if err != nil {
			return nil, errors.Wrap(err, "error on scan row")
		}

		entity := m.ToEntity()
		res = append(res, &entity)
	}

	return res, nil
}

func (repo *postRepo) ListPublished(ctx context.Context) ([]*post.Entity, error) {
	query := `SELECT uuid, title, slug, content_markdown, content_html, published_at FROM posts WHERE published_at IS NOT NULL;`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on query")
	}

	res := make([]*post.Entity, 0)
	for rows.Next() {
		var m postModel
		dest := []interface{}{&m.UUID, &m.Title, &m.Slug, &m.ContentMarkdown, &m.ContentHTML, &m.PublishedAt}
		err = rows.Scan(dest...)
		if err != nil {
			return nil, errors.Wrap(err, "error on scan row")
		}

		entity := m.ToEntity()
		res = append(res, &entity)
	}

	return res, nil
}

func (repo *postRepo) FindByUUID(ctx context.Context, uuid string) (*post.Entity, error) {
	var m postModel

	// prepare query
	query := `SELECT uuid, title, slug, content_markdown, content_html, published_at FROM posts WHERE uuid = $1;`
	args := []interface{}{uuid}
	dest := []interface{}{&m.UUID, &m.Title, &m.Slug, &m.ContentMarkdown, &m.ContentHTML, &m.PublishedAt}

	err := repo.db.QueryRowContext(ctx, query, args...).Scan(dest...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(errors.WithStack(err), "error on query row")
	}

	entity := m.ToEntity()

	return &entity, nil
}

func (repo *postRepo) FindBySlug(ctx context.Context, slug string) (*post.Entity, error) {
	var m postModel

	// prepare query
	query := `SELECT uuid, title, slug, content_markdown, content_html, published_at FROM posts WHERE slug = $1;`
	args := []interface{}{slug}
	dest := []interface{}{&m.UUID, &m.Title, &m.Slug, &m.ContentMarkdown, &m.ContentHTML, &m.PublishedAt}

	err := repo.db.QueryRowContext(ctx, query, args...).Scan(dest...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(errors.WithStack(err), "error on query row")
	}

	entity := m.ToEntity()

	return &entity, nil
}

func (repo *postRepo) Insert(ctx context.Context, entity post.Entity) error {
	m := new(postModel)
	m.FromEntity(entity)

	query := `INSERT INTO posts (uuid, title, slug, content_markdown, content_html, published_at) VALUES ($1, $2, $3, $4, $5, $6);`
	args := []interface{}{m.UUID, m.Title, m.Slug, m.ContentMarkdown, m.ContentHTML, m.PublishedAt}

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
