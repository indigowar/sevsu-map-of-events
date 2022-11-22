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
	organizer :=
		models.Organizer{ID: uuid.New(), Name: name, Logo: logo, Level: level}
	err := o.storage.Create(ctx, organizer)
	if err != nil {
		return models.Organizer{}, err
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
	level := models.OrganizerLevel{ID: uuid.New(), Name: name, Code: code}

	if err := o.storage.AddLevel(ctx, level); err != nil {
		return models.OrganizerLevel{}, err
	}
	return level, nil
}

func (o organizerSvc) UpdateLevel(_ context.Context, _ models.OrganizerLevel) (models.OrganizerLevel, error) {
	panic("unimplemented")
}

func (o organizerSvc) Update(ctx context.Context, id uuid.UUID, name, logo string, level uuid.UUID) (models.Organizer, error) {
	m := models.Organizer{
		ID:    id,
		Name:  name,
		Logo:  logo,
		Level: level,
	}

	if err := o.storage.Update(ctx, m); err != nil {
		log.Println(err)
		return models.Organizer{}, errors.New("failed to update an organizer")
	}

	return m, nil
}

func NewOrganizerService(storage storages.OrganizerStorageRepository) (services.OrganizerService, error) {
	return &organizerSvc{
		storage: storage,
	}, nil
}
