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

type PostgresCompetitorStorage struct {
	con *pgx.Conn
}

func (storage PostgresCompetitorStorage) Get(ctx context.Context, id uuid.UUID) (models.Competitor, error) {
	var Id uuid.UUID
	var name string

	query := fmt.Sprintf("SELECT * FROM competitor WHERE id == '%s'", id.String())

	if err := storage.con.QueryRow(ctx, query).Scan(&Id, &name); err != nil {
		log.Println("Got query error or scan error: ", err)
		return nil, err
	}

	return models.NewCompetitor(Id, name), nil
}

func (storage PostgresCompetitorStorage) GetAll(ctx context.Context) ([]models.Competitor, error) {
	//TODO implement me
	panic("implement me")
}

func (storage PostgresCompetitorStorage) Create(ctx context.Context, competitor models.Competitor) error {
	//TODO implement me
	panic("implement me")
}

func (storage PostgresCompetitorStorage) Update(ctx context.Context, competitor models.Competitor) error {
	//TODO implement me
	panic("implement me")
}

func (storage PostgresCompetitorStorage) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresCompetitorStorage(con *pgx.Conn) (storages.CompetitorStorage, error) {
	if con == nil {
		return nil, errors.New("invalid connection")
	}

	if con.IsClosed() {
		return nil, errors.New("connection is closed")
	}

	return &PostgresCompetitorStorage{con: con}, nil
}
