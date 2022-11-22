package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type CompetitorStorage interface {
	AllIDs(ctx context.Context) ([]uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (models.Competitor, error)
	GetAll(ctx context.Context) ([]models.Competitor, error)
	Create(ctx context.Context, competitor models.Competitor) error
	Update(ctx context.Context, competitor models.Competitor) error
	Delete(ctx context.Context, id uuid.UUID) error

	StorageTransaction
}
