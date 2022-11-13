package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type EventStorageRepository interface {
	GetIDList(ctx context.Context) ([]uuid.UUID, error)
	Filter(ctx context.Context) ([]uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Event, error)
	Create(ctx context.Context, event models.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, event models.Event) error
}
