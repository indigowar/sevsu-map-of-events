package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type competitorService struct {
	storage storages.CompetitorStorage
}

func (svc competitorService) GetByID(ctx context.Context, id uuid.UUID) (models.Competitor, error) {
	c, err := svc.storage.Get(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("not found")
	}
	return c, nil
}

func (svc competitorService) GetAll(ctx context.Context) ([]models.Competitor, error) {
	competitors, err := svc.storage.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return competitors, nil
}

func (svc competitorService) Create(ctx context.Context, name string) (models.Competitor, error) {
	c := models.NewCompetitor(uuid.New(), name)

	if err := svc.storage.Create(ctx, c); err != nil {
		log.Println(err)
		return nil, errors.New("failed to create a competitor")
	}

	return c, nil
}

func (svc competitorService) Delete(ctx context.Context, id uuid.UUID) error {
	return svc.storage.Delete(ctx, id)
}

func NewCompetitorService(storage storages.CompetitorStorage) services.CompetitorService {
	return &competitorService{
		storage: storage,
	}
}
