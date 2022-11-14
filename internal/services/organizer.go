package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type organizerSvc struct {
	storage storages.OrganizerStorageRepository
}

func (o organizerSvc) GetAll(ctx context.Context) ([]models.Organizer, error) {
	return o.storage.GetAll(ctx)
}

func (o organizerSvc) GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error) {
	return o.storage.GetByID(ctx, id)
}

func (o organizerSvc) Create(ctx context.Context, name, logo string, level uuid.UUID) (models.Organizer, error) {
	// TODO: add validation of user input
	organizer := models.NewOrganizer(uuid.New(), name, logo, level)
	err := o.storage.Create(ctx, organizer)
	if err != nil {
		return nil, err
	}
	return organizer, nil
}

func (o organizerSvc) Delete(ctx context.Context, id uuid.UUID) error {
	return o.storage.Delete(ctx, id)
}

func (o organizerSvc) GetAllLevels(ctx context.Context) ([]models.OrganizerLevel, error) {
	return o.storage.GetLevels(ctx)
}

func (o organizerSvc) CreateLevel(ctx context.Context, name string, code string) (models.OrganizerLevel, error) {
	// TODO: add validation
	level := models.NewOrganizerLevel(uuid.New(), name, code)

	if err := o.storage.AddLevel(ctx, level); err != nil {
		return nil, err
	}
	return level, nil
}

func (o organizerSvc) UpdateLevel(ctx context.Context, level models.OrganizerLevel) (models.OrganizerLevel, error) {
	panic("unimplemented")
}

func NewOrganizerService(storage storages.OrganizerStorageRepository) (services.OrganizerService, error) {
	return &organizerSvc{
		storage: storage,
	}, nil
}
