package metadata

import (
	"context"
	"errors"

	"movieexample.com/src/metadata/internal/repository"
	"movieexample.com/src/metadata/pkg/model"
)

var ErrNotFound = errors.New("not Found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

type Controller struct {
	repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, nil
}
