package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type RangeStorageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.RangeModel, error)
	GetMaximumRange(ctx context.Context) (models.RangeModel, error)
	Create(ctx context.Context, foundingRange models.RangeModel) (models.RangeModel, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
