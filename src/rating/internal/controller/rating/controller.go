package rating

import (
	"context"
	"errors"

	"movieexample.com/src/rating/internal/repository"
	"movieexample.com/src/rating/pkg/model"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines a arating service controller
type Controller struct {
	repo     ratingRepository
	ingester ratingIngester
}

type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

// New creates a rating service controller
func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo, ingester}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it
func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && err == repository.ErrNotFound {
		return 0, ErrNotFound
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

func (c *Controller) StartIngestion(ctx context.Context) error {
	ch, err := c.ingester.Ingest(ctx)
	if err != nil {
		return err
	}
	for e := range ch {
		if err := c.PutRating(ctx, e.RecordID, e.RecordType, &model.Rating{
			UserID: e.UserID,
			Value:  e.Value,
		}); err != nil {
			return err
		}
	}
	return nil
}

// PutRating writes a rating for a given record
func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}
