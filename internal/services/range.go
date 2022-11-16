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
	storage        storages.RangeStorageRepository
	validatorsList []func(foundingRange models.RangeModel) error
}

func (svc rangeService) GetByID(ctx context.Context, id uuid.UUID) (models.RangeModel, error) {
	r, err := svc.storage.GetByID(ctx, id)
	if err != nil {
		log.Println(err)
		return models.RangeModel{}, errors.New("internal error")
	}

	return r, nil
}

func (svc rangeService) GetMaximumRange(ctx context.Context) (models.RangeModel, error) {
	result, err := svc.storage.GetMaximumRange(ctx)
	if err != nil {
		log.Println(err)
		return models.RangeModel{}, errors.New("internal error")
	}
	return result, nil
}

func (svc rangeService) Create(ctx context.Context, low, high int) (models.RangeModel, error) {
	r := models.RangeModel{ID: uuid.New(), Low: low, High: high}

	for _, validator := range svc.validatorsList {
		if err := validator(r); err != nil {
			return models.RangeModel{}, err
		}
	}

	return svc.storage.Create(ctx, r)
}

func (svc rangeService) Delete(ctx context.Context, id uuid.UUID) error {
	return svc.Delete(ctx, id)
}

func (svc rangeService) Update(ctx context.Context, foundingRange models.RangeModel) (models.RangeModel, error) {
	for _, validator := range svc.validatorsList {
		if err := validator(foundingRange); err != nil {
			return models.RangeModel{}, err
		}
	}

	err := svc.Delete(ctx, foundingRange.ID)
	if err != nil {
		return models.RangeModel{}, err
	}
	return svc.storage.Create(ctx, foundingRange)
}

func NewFoundingRangeService(storage storages.RangeStorageRepository) services.RangeService {
	return &rangeService{
		storage:        storage,
		validatorsList: []func(foundingRange models.RangeModel) error{validators.ValidateRange},
	}
}

func NewCoFoundingRangeService(storage storages.RangeStorageRepository) services.RangeService {
	return &rangeService{
		storage: storage,
		validatorsList: []func(foundingRange models.RangeModel) error{
			validators.ValidateRange,
			validators.ValidatePercentRange,
		},
	}
}
