package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type CompetitorService interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.Competitor, error)
	GetAll(ctx context.Context) ([]models.Competitor, error)
	Create(ctx context.Context, name string) (models.Competitor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
