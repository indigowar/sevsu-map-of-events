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

func (s PostgresCompetitorStorage) Get(ctx context.Context, id uuid.UUID) (models.Competitor, error) {
	var Id uuid.UUID
	var name string

	query := fmt.Sprintf("SELECT * FROM competitor WHERE competitor_id == '%s'", id.String())

	if err := s.con.QueryRow(ctx, query).Scan(&Id, &name); err != nil {
		log.Println("Got query error or scan error: ", err)
		return nil, err
	}

	return models.NewCompetitor(Id, name), nil
}

func (s PostgresCompetitorStorage) GetAll(ctx context.Context) ([]models.Competitor, error) {
	comps := make([]models.Competitor, 0)

	rows, err := s.con.Query(ctx, "SELECT * FROM competitor")
	if err != nil {
		log.Println("Failed to read from database")
		return nil, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println("Failed to read fetched value from database")
			return nil, err
		}

		id := values[0].([16]byte)
		name := values[1].(string)

		parsedId, _ := uuid.FromBytes(id[:])

		comps = append(comps, models.NewCompetitor(parsedId, name))
	}

	return comps, nil
}

func (s PostgresCompetitorStorage) Create(ctx context.Context, competitor models.Competitor) error {
	command := "INSERT INTO competitor (competitor_id, competitor_name) VALUES ($1, $2)"
	if _, err := s.con.Exec(ctx, command, competitor.ID(), competitor.Name()); err != nil {
		log.Println(err)
		return errors.New("failed to create new competitor")
	}
	return nil
}

func (s PostgresCompetitorStorage) Update(ctx context.Context, competitor models.Competitor) error {
	command := "UPDATE competitor SET competitor_name = $2 WHERE competitor_id = $1"
	if _, err := s.con.Exec(ctx, command, competitor.ID(), competitor.Name()); err != nil {
		log.Println(err)
		return errors.New("failed to update a competitor")
	}
	return nil
}

func (s PostgresCompetitorStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.con.Exec(ctx, "DELETE FROM competitor WHERE competitor_id = $1", id)
	return err
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
