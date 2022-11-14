package storages

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
)

type postgresEventStorage struct {
	con *pgx.Conn
}

func (s postgresEventStorage) GetIDList(ctx context.Context) ([]uuid.UUID, error) {
	// query := "SELECT * FROM event JOIN competitor_requirements cr on event.id = cr.event JOIN subject s on event.id = s.event"
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) Filter(ctx context.Context) ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) Create(ctx context.Context, event models.Event) error {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) Update(ctx context.Context, event models.Event) error {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) AddSubject(ctx context.Context, id, subjectId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) RemoveSubject(ctx context.Context, id, subjectId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) AddCompetitor(ctx context.Context, id, competitorId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresEventStorage(con *pgx.Conn) storages.EventStorageRepository {
	return &postgresEventStorage{
		con: con,
	}
}
