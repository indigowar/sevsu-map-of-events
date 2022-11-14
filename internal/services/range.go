package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/internal/domain/services"
	"github.com/indigowar/map-of-events/internal/domain/validators"
)

type rangeService struct {
	storage storages.RangeStorageRepository
}

func (svc rangeService) GetByID(ctx context.Context, id uuid.UUID) (models.FoundingRange, error) {
	r, err := svc.storage.GetByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	return r, nil
}

func (svc rangeService) GetMaximumRange(ctx context.Context) (models.FoundingRange, error) {
	result, err := svc.storage.GetMaximumRange(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return result, nil
}

func (svc rangeService) Create(ctx context.Context, low, high int) (models.FoundingRange, error) {
	r := models.NewRange(uuid.New(), low, high)

	if err := validators.ValidateRange(r); err != nil {
		return nil, err
	}

	return svc.storage.Create(ctx, r)
}

func (svc rangeService) Delete(ctx context.Context, id uuid.UUID) error {
	return svc.Delete(ctx, id)
}

func (svc rangeService) Update(ctx context.Context, foundingRange models.FoundingRange) (models.FoundingRange, error) {
	if err := validators.ValidateRange(foundingRange); err != nil {
		return nil, err
	}

	err := svc.Delete(ctx, foundingRange.ID())
	if err != nil {
		return nil, err
	}
	return svc.storage.Create(ctx, foundingRange)
}

func NewFoundingRangeService(storage storages.RangeStorageRepository) services.RangeService {
	return &rangeService{
		storage: storage,
	}
}
