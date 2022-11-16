package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type eventService struct {
	organizer        services.OrganizerService
	foundingRanges   services.RangeService
	coFoundingRanges services.RangeService
	competitors      services.CompetitorService

	eventStorage storages.EventStorageRepository
}

func (svc eventService) GetAll(ctx context.Context) ([]models.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (svc eventService) GetByID(ctx context.Context, id uuid.UUID) (models.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (svc eventService) Create(ctx context.Context, event models.Event) (models.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (svc eventService) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (svc eventService) GetAllAsMinimal(ctx context.Context, id uuid.UUID) ([]services.EventMinimal, error) {
	//TODO implement me
	panic("implement me")
}

func (svc eventService) GetByIDAsMinimal(ctx context.Context, id uuid.UUID) (services.EventMinimal, error) {
	//TODO implement me
	panic("implement me")
}

func NewEventServices(storage storages.EventStorageRepository,
	organizer services.OrganizerService,
	foundingRange, coFoundingRange services.RangeService,
	competitors services.CompetitorService) services.EventService {
	return &eventService{
		organizer:        organizer,
		foundingRanges:   foundingRange,
		coFoundingRanges: coFoundingRange,
		competitors:      competitors,
		eventStorage:     storage,
	}
}
