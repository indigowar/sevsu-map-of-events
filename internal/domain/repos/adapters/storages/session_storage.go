package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type SessionStorage interface {
	GetByToken(ctx context.Context, token string) (models.TokenSession, error)
	GetByUser(ctx context.Context, id uuid.UUID) (models.TokenSession, error)

	Create(ctx context.Context, session models.TokenSession) error
	Delete(ctx context.Context, token string) error
}
