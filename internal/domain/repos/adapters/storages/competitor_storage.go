package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/pkg/errors"
)

type CompetitorStorage interface {
	AllIDs(ctx context.Context) ([]uuid.UUID, errors.Error)
	Get(ctx context.Context, id uuid.UUID) (models.Competitor, errors.Error)
	GetAll(ctx context.Context) ([]models.Competitor, errors.Error)
	Create(ctx context.Context, competitor models.Competitor) errors.Error
	Update(ctx context.Context, competitor models.Competitor) errors.Error
	Delete(ctx context.Context, id uuid.UUID) errors.Error

	StorageTransaction
}
