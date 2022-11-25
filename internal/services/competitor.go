package services

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/services"
	"github.com/indigowar/map-of-events/pkg/errors"
)

type competitorService struct {
	storage adapters.CompetitorStorage
}

func (svc competitorService) AllIDs(ctx context.Context) ([]uuid.UUID, errors.Error) {
	result, err := svc.storage.AllIDs(ctx)
	if err != nil {
		return nil, failedToReadDatabaseErr(err, "get all IDs")
	}
	return result, nil
}

func (svc competitorService) GetByID(ctx context.Context, id uuid.UUID) (models.Competitor, errors.Error) {
	c, err := svc.storage.Get(ctx, id)
	if err != nil {
		log.Println(err)
		return models.Competitor{}, failedToReadDatabaseErr(err, "get competitor by id")
	}
	return c, nil
}

func (svc competitorService) GetAll(ctx context.Context) ([]models.Competitor, errors.Error) {
	competitors, err := svc.storage.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return nil, failedToReadDatabaseErr(err, "get all competitors")
	}
	return competitors, nil
}

func (svc competitorService) Create(ctx context.Context, name string) (models.Competitor, errors.Error) {
	c := models.Competitor{
		ID:   uuid.New(),
		Name: name,
	}

	if err := svc.storage.Create(ctx, c); err != nil {
		log.Println(err)
		return models.Competitor{}, failedToReadDatabaseErr(err, "create a competitor")
	}
	return c, nil
}

func (svc competitorService) Delete(ctx context.Context, id uuid.UUID) errors.Error {
	if err := svc.storage.Delete(ctx, id); err != nil {
		log.Println(err)
		return failedToReadDatabaseErr(err, "delete a competitor")
	}
	return nil
}

func NewCompetitorService(storage adapters.CompetitorStorage) services.CompetitorService {
	return &competitorService{
		storage: storage,
	}
}

// failedToReadDatabaseErr - creates an error from e error that explains why failed to read the database
func failedToReadDatabaseErr(e errors.Error, targetOfJob string) errors.Error {
	return errors.CreateError(services.ErrReasonInternalError,
		fmt.Sprintf("failed to %s: internal server error", targetOfJob),
		fmt.Sprintf("failed to %s: internal server error, because: %s", targetOfJob, e.ShortErr()),
	)
}
