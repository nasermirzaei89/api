package post

import (
	"context"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type service struct {
	repo Repository
}

func (svc *service) UpdatePostByUUID(ctx context.Context, postUUID string, req UpdatePostByUUIDRequest) (*Entity, error) {
	entity, err := svc.repo.FindByUUID(ctx, postUUID)
	if err != nil {
		return nil, errors.Wrap(err, "error on find post by uuid")
	}

	if entity == nil {
		return nil, ErrPostWithUUIDNotFound{UUID: postUUID}
	}

	if req.Slug == "" {
		req.Slug = req.Title
	}

	req.Slug = slug.Make(req.Slug)

	index := 1
	uniqueSlug := req.Slug
	for {
		entity, err := svc.repo.FindBySlug(ctx, uniqueSlug)
		if err != nil {
			return nil, errors.Wrap(err, "error on find post by slug")
		}

		if entity == nil || entity.UUID == postUUID {
			req.Slug = uniqueSlug
			break
		}

		index++
		uniqueSlug = fmt.Sprintf("%s-%d", req.Slug, index)
	}

	contentHTML := string(markdown.ToHTML([]byte(req.ContentMarkdown), parser.New(), nil))

	entity.Title = req.Title
	entity.Slug = req.Slug
	entity.ContentMarkdown = req.ContentMarkdown
	entity.ContentHTML = contentHTML

	err = svc.repo.UpdateByUUID(ctx, postUUID, *entity)
	if err != nil {
		return nil, errors.Wrap(err, "error on update post by uuid")
	}

	return entity, nil
}

func (svc *service) GetPostByUUID(ctx context.Context, postUUID string) (*Entity, error) {
	entity, err := svc.repo.FindByUUID(ctx, postUUID)
	if err != nil {
		return nil, errors.Wrap(err, "error on find post by uuid")
	}

	if entity == nil {
		return nil, ErrPostWithUUIDNotFound{UUID: postUUID}
	}

	return entity, nil
}

func (svc *service) GetPostBySlug(ctx context.Context, slug string) (*Entity, error) {
	entity, err := svc.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.Wrap(err, "error on find post by slug")
	}

	if entity == nil {
		return nil, ErrPostWithSlugNotFound{Slug: slug}
	}

	return entity, nil
}

func (svc *service) ListPosts(ctx context.Context) ([]*Entity, error) {
	res, err := svc.repo.List(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error on list posts")
	}

	return res, nil
}

func (svc *service) CreatePost(ctx context.Context, req CreatePostRequest) (*Entity, error) {
	if req.Slug == "" {
		req.Slug = req.Title
	}

	req.Slug = slug.Make(req.Slug)

	index := 1
	uniqueSlug := req.Slug
	for {
		entity, err := svc.repo.FindBySlug(ctx, uniqueSlug)
		if err != nil {
			return nil, errors.Wrap(err, "error on find post by slug")
		}

		if entity == nil {
			req.Slug = uniqueSlug
			break
		}

		index++
		uniqueSlug = fmt.Sprintf("%s-%d", req.Slug, index)
	}

	contentHTML := string(markdown.ToHTML([]byte(req.ContentMarkdown), parser.New(), nil))

	entity := Entity{
		UUID:            uuid.New().String(),
		Title:           req.Title,
		Slug:            req.Slug,
		ContentMarkdown: req.ContentMarkdown,
		ContentHTML:     contentHTML,
	}

	err := svc.repo.Insert(ctx, entity)
	if err != nil {
		return nil, errors.Wrap(err, "error on insert post")
	}

	return &entity, nil
}

func NewService(repo Repository) Service {
	svc := service{repo: repo}

	return &svc
}
