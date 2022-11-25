package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/pkg/postgres"
)

type PostgresOrganizerStorage struct {
	pool *pgxpool.Pool
}

func (s PostgresOrganizerStorage) GetAllIDs(ctx context.Context) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)

	rows, err := s.pool.Query(ctx, "SELECT organizer_id FROM organizer")
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read database")
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read fetched values")
		}
		parsedId := values[0].([16]byte)
		id, err := uuid.FromBytes(parsedId[:])
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to parse id")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (s PostgresOrganizerStorage) GetLevelsIDs(ctx context.Context) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)

	rows, err := s.pool.Query(ctx, "SELECT organizer_level_id FROM organizer_level")
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read database")
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read fetched values")
		}
		parsedId := values[0].([16]byte)
		id, err := uuid.FromBytes(parsedId[:])
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to parse id")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (s PostgresOrganizerStorage) BeginTransaction(ctx context.Context) (interface{}, error) {
	return s.pool.Begin(ctx)
}

func (s PostgresOrganizerStorage) CloseTransaction(ctx context.Context, transaction interface{}) error {
	tx := transaction.(*pgxpool.Tx)
	return tx.Rollback(ctx)
}

func NewPostgresOrganizerStorage(p *pgxpool.Pool) adapters.OrganizerStorage {
	return &PostgresOrganizerStorage{
		pool: p,
	}
}

func (s PostgresOrganizerStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	var organizer models.Organizer

	command := fmt.Sprintf("SELECT * FROM organizer WHERE organizer_id = '%s'", id.String())

	err := dataSource.QueryRow(ctx, command).Scan(&organizer.ID, &organizer.Name, &organizer.Logo, &organizer.Level)

	if err != nil {
		log.Println("Failed to read to database")
		return models.Organizer{}, err
	}

	return organizer, nil
}

func (s PostgresOrganizerStorage) GetAll(ctx context.Context) ([]models.Organizer, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	organizers := make([]models.Organizer, 0)

	rows, err := dataSource.Query(ctx, "SELECT * FROM organizer")
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to query database")
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read fetched values")
		}
		parsedId := values[0].([16]byte)
		id, err := uuid.FromBytes(parsedId[:])
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to parse id")
		}
		name := values[1].(string)
		logo := values[2].(string)
		byteLevel := values[3].([16]byte)
		level, _ := uuid.FromBytes(byteLevel[:])

		organizers = append(organizers, models.Organizer{ID: id, Name: name, Logo: logo, Level: level})
	}

	return organizers, nil
}

func (s PostgresOrganizerStorage) Add(ctx context.Context, organizer models.Organizer) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "INSERT INTO organizer (organizer_id, organizer_name, organizer_image, organizer_level) VALUES ($1, $2, $3, $4)"

	if _, err := dataSource.Exec(ctx, command, organizer.ID, organizer.Name, organizer.Logo, organizer.Level); err != nil {
		log.Println(err)
		return errors.New("failed to create organizer")
	}
	return nil
}

func (s PostgresOrganizerStorage) Remove(ctx context.Context, id uuid.UUID) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	_, err := dataSource.Exec(ctx, "DELETE FROM organizer WHERE organizer_id=$1", id)
	return err
}

func (s PostgresOrganizerStorage) Update(ctx context.Context, organizer models.Organizer) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "UPDATE organizer SET organizer_name = $2, organizer_image = $3, organizer_level = $4 WHERE organizer_id = $1"

	if _, err := dataSource.Exec(ctx, command, organizer.ID, organizer.Name, organizer.Logo, organizer.Level); err != nil {
		log.Println(err)
		return errors.New("failed to update")
	}

	return nil
}

func (s PostgresOrganizerStorage) GetLevels(ctx context.Context) ([]models.OrganizerLevel, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	levels := make([]models.OrganizerLevel, 0)

	query := "SELECT * FROM organizer_level"
	rows, err := dataSource.Query(ctx, query)

	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read database")
	}

	for rows.Next() {
		values, err := rows.Values()

		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read fetched values")
		}

		idInBytes := values[0].([16]byte)

		id, err := uuid.FromBytes(idInBytes[:])
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read parse id from database")
		}

		name := values[1].(string)
		code := values[2].(string)
		levels = append(levels, models.OrganizerLevel{ID: id, Name: name, Code: code})
	}

	return levels, nil
}

func (s PostgresOrganizerStorage) AddLevel(ctx context.Context, level models.OrganizerLevel) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "INSERT INTO organizer_level(organizer_level_id, organizer_level_name, organizer_level_code) VALUES ($1, $2, $3)"

	if _, err := dataSource.Exec(ctx, command, level.ID, level.Name, level.Code); err != nil {
		log.Println(err)
		return errors.New("failed to add to database")
	}
	return nil
}

func (s PostgresOrganizerStorage) RemoveLevel(ctx context.Context, id uuid.UUID) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "DELETE FROM organizer_level WHERE organizer_level_id = $1"

	if _, err := dataSource.Exec(ctx, command, id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}
	return nil
}
