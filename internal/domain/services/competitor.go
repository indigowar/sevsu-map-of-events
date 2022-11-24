package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/pkg/errors"
)

type CompetitorService interface {
	AllIDs(ctx context.Context) ([]uuid.UUID, errors.Error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Competitor, errors.Error)
	GetAll(ctx context.Context) ([]models.Competitor, errors.Error)
	Create(ctx context.Context, name string) (models.Competitor, errors.Error)
	Delete(ctx context.Context, id uuid.UUID) errors.Error
}
