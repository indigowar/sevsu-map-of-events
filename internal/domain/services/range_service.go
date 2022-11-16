package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type RangeService interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.RangeModel, error)
	GetMaximumRange(ctx context.Context) (models.RangeModel, error)
	Create(ctx context.Context, low, high int) (models.RangeModel, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, foundingRange models.RangeModel) (models.RangeModel, error)
}
