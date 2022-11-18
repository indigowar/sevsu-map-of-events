package storages

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/pkg/postgres"
)

type PostgresCompetitorStorage struct {
	pool *pgxpool.Pool
}

func (s PostgresCompetitorStorage) Get(ctx context.Context, id uuid.UUID) (models.Competitor, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	var Id uuid.UUID
	var name string

	query := fmt.Sprintf("SELECT * FROM competitor WHERE competitor_id == '%s'", id.String())

	if err := dataSource.QueryRow(ctx, query).Scan(&Id, &name); err != nil {
		log.Println("Got query error or scan error: ", err)
		return models.Competitor{}, err
	}

	return models.Competitor{ID: id, Name: name}, nil
}

func (s PostgresCompetitorStorage) GetAll(ctx context.Context) ([]models.Competitor, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	comps := make([]models.Competitor, 0)

	rows, err := dataSource.Query(ctx, "SELECT * FROM competitor")
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

		comps = append(comps, models.Competitor{ID: parsedId, Name: name})
	}

	return comps, nil
}

func (s PostgresCompetitorStorage) Create(ctx context.Context, competitor models.Competitor) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "INSERT INTO competitor (competitor_id, competitor_name) VALUES ($1, $2)"
	if _, err := dataSource.Exec(ctx, command, competitor.ID, competitor.Name); err != nil {
		log.Println(err)
		return errors.New("failed to create new competitor")
	}
	return nil
}

func (s PostgresCompetitorStorage) Update(ctx context.Context, competitor models.Competitor) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "UPDATE competitor SET competitor_name = $2 WHERE competitor_id = $1"
	if _, err := dataSource.Exec(ctx, command, competitor.ID, competitor.Name); err != nil {
		log.Println(err)
		return errors.New("failed to update a competitor")
	}
	return nil
}

func (s PostgresCompetitorStorage) Delete(ctx context.Context, id uuid.UUID) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	_, err := dataSource.Exec(ctx, "DELETE FROM competitor WHERE competitor_id = $1", id)
	return err
}

func NewPostgresCompetitorStorage(p *pgxpool.Pool) storages.CompetitorStorage {
	return &PostgresCompetitorStorage{pool: p}
}
