package services

import (
	"context"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type ImageService interface {
	GetAllLinks(ctx context.Context) ([]string, error)
	Get(ctx context.Context, link string) (models.StoredImage, error)
	Create(ctx context.Context, link string, image []byte) (models.StoredImage, error)
	Delete(ctx context.Context, link string) error
	Update(ctx context.Context, link string, image []byte) (models.StoredImage, error)
}
