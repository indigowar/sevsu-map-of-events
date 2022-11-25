package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type UserStorage interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetByName(ctx context.Context, name string) (models.User, error)
	Create(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateName(ctx context.Context, id uuid.UUID, name string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
}
