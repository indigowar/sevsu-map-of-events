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
)

type PostgresOrganizerStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresOrganizerStorage(p *pgxpool.Pool) storages.OrganizerStorageRepository {
	return &PostgresOrganizerStorage{
		pool: p,
	}
}

func (s PostgresOrganizerStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error) {
	var Id, level uuid.UUID
	var name, logo string
	err := s.pool.QueryRow(ctx, fmt.Sprintf("SELECT * FROM organizer WHERE organizer_id = '%s'", id.String())).Scan(&Id, &name, &logo, &level)
	if err != nil {
		log.Println("Failed to read to database")
		return models.Organizer{}, err
	}

	return models.Organizer{ID: Id, Name: name, Logo: logo, Level: level}, nil
}

func (s PostgresOrganizerStorage) GetAll(ctx context.Context) ([]models.Organizer, error) {
	organizers := make([]models.Organizer, 0)

	rows, err := s.pool.Query(ctx, "SELECT * FROM organizer")
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
		id := values[0].(uuid.UUID)
		name := values[1].(string)
		logo := values[2].(string)
		level := values[3].(uuid.UUID)

		organizers = append(organizers, models.Organizer{ID: id, Name: name, Logo: logo, Level: level})
	}

	return organizers, nil
}

func (s PostgresOrganizerStorage) Create(ctx context.Context, organizer models.Organizer) error {
	command := "INSERT INTO organizer (organizer_id, organizer_name, organizer_image, organizer_level) VALUES ($1, $2, $3, $4)"

	if _, err := s.pool.Exec(ctx, command, organizer.ID, organizer.Name, organizer.Logo, organizer.Level); err != nil {
		log.Println(err)
		return errors.New("failed to create organizer")
	}
	return nil
}

func (s PostgresOrganizerStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM organizer WHERE organizer_id=$1", id)
	return err
}

func (s PostgresOrganizerStorage) Update(ctx context.Context, organizer models.Organizer) error {
	command := "UPDATE organizer SET organizer_name = $2, organizer_image = $3, organizer_level = $4 WHERE organizer_id = $1"

	if _, err := s.pool.Exec(ctx, command, organizer.ID, organizer.Name, organizer.Logo, organizer.Level); err != nil {
		log.Println(err)
		return errors.New("failed to update")
	}

	return nil
}

func (s PostgresOrganizerStorage) GetLevels(ctx context.Context) ([]models.OrganizerLevel, error) {
	levels := make([]models.OrganizerLevel, 0)

	query := "SELECT * FROM organizer_level"
	rows, err := s.pool.Query(ctx, query)

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

		id := values[0].(uuid.UUID)
		name := values[1].(string)
		code := values[2].(string)
		levels = append(levels, models.OrganizerLevel{ID: id, Name: name, Code: code})
	}

	return levels, nil
}

func (s PostgresOrganizerStorage) AddLevel(ctx context.Context, level models.OrganizerLevel) error {
	command := "INSERT INTO organizer_level(organizer_level_id, organizer_level_name, organizer_level_code) VALUES ($1, $2, $3)"

	if _, err := s.pool.Exec(ctx, command, level.ID, level.Name, level.Code); err != nil {
		log.Println(err)
		return errors.New("failed to add to database")
	}
	return nil
}

func (s PostgresOrganizerStorage) DeleteLevel(ctx context.Context, id uuid.UUID) error {
	command := "DELETE FROM organizer_level WHERE organizer_level_id = $1"

	if _, err := s.pool.Exec(ctx, command, id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}
	return nil
}
