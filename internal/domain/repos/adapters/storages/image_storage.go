package storages

import (
	"context"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type ImageStorage interface {
	Get(ctx context.Context, link string) (models.StoredImage, error)
	Create(ctx context.Context, image models.StoredImage) error
	Delete(ctx context.Context, link string) error
	Update(ctx context.Context, image models.StoredImage) error
}
