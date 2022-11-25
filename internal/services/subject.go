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

type subjectService struct {
	storage storages.SubjectStorage
}

func (svc subjectService) GetAllExisting(ctx context.Context) ([]models.Subject, error) {
	return svc.storage.GetAll(ctx)
}

func (svc subjectService) GetAllForEvent(ctx context.Context, eventId uuid.UUID) ([]models.Subject, error) {
	ids, err := svc.storage.GetByEvent(ctx, eventId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := make([]models.Subject, len(ids))

	for i, v := range ids {
		result[i], _ = svc.storage.GetByID(ctx, v)
	}

	return result, nil
}

func (svc subjectService) GetByID(ctx context.Context, id uuid.UUID) (models.Subject, error) {
	return svc.storage.GetByID(ctx, id)
}

func (svc subjectService) Create(ctx context.Context, eventId uuid.UUID, subject string) (models.Subject, error) {
	s := models.Subject{
		ID:      uuid.New(),
		Name:    subject,
		EventID: eventId,
	}

	if err := svc.storage.Add(ctx, s); err != nil {
		log.Println(err)
		return models.Subject{}, errors.New("failed to add")
	}

	return s, nil
}

func (svc subjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return svc.storage.Delete(ctx, id)
}

func (svc subjectService) Update(ctx context.Context, subject models.Subject) error {
	return svc.Update(ctx, subject)
}

func NewSubjectService(storage storages.SubjectStorage) services.SubjectService {
	return &subjectService{
		storage: storage,
	}
}
