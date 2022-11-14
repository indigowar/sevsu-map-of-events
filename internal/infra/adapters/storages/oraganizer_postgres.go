package storages

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
)

type PostgresOrganizerStorage struct {
	conn *pgx.Conn
}

func NewPostgresOrganizerStorage(conn *pgx.Conn) (storages.OrganizerStorageRepository, error) {
	if conn.IsClosed() || conn == nil {
		return nil, errors.New("invalid connection")
	}

	return &PostgresOrganizerStorage{
		conn: conn,
	}, nil
}

func (p PostgresOrganizerStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error) {
	var Id, level uuid.UUID
	var name, logo string
	err := p.conn.QueryRow(ctx, fmt.Sprintf("SELECT * FROM organizer WHERE id == '%s'", id.String())).Scan(&Id, &name, &logo, &level)
	if err != nil {
		log.Println("Failed to read to database")
		return nil, err
	}

	return models.NewOrganizer(Id, name, logo, level), nil
}

func (p PostgresOrganizerStorage) GetAll(ctx context.Context) ([]models.Organizer, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrganizerStorage) Create(ctx context.Context, organizer models.Organizer) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrganizerStorage) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrganizerStorage) Update(ctx context.Context, organizer models.Organizer) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrganizerStorage) GetLevels(ctx context.Context) ([]models.OrganizerLevel, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrganizerStorage) AddLevel(ctx context.Context, level models.OrganizerLevel) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrganizerStorage) DeleteLevel(ctx context.Context, id uuid.UUID, replacement uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
