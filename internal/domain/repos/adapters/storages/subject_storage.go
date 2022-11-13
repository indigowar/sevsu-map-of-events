package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type SubjectStorageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.Subject, error)
}
