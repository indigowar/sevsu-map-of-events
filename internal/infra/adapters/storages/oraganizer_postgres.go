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

func (s PostgresOrganizerStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error) {
	var Id, level uuid.UUID
	var name, logo string
	err := s.conn.QueryRow(ctx, fmt.Sprintf("SELECT * FROM organizer WHERE id == '%s'", id.String())).Scan(&Id, &name, &logo, &level)
	if err != nil {
		log.Println("Failed to read to database")
		return nil, err
	}

	return models.NewOrganizer(Id, name, logo, level), nil
}

func (s PostgresOrganizerStorage) GetAll(ctx context.Context) ([]models.Organizer, error) {
	organizers := make([]models.Organizer, 0)

	rows, err := s.conn.Query(ctx, "SELECT * FROM organizer")
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

		organizers = append(organizers, models.NewOrganizer(id, name, logo, level))
	}

	return organizers, nil
}

func (s PostgresOrganizerStorage) Create(ctx context.Context, organizer models.Organizer) error {
	command := "INSERT INTO organizer (id, name, logo, level) VALUES ($1, $2, $3, $4)"

	if _, err := s.conn.Exec(ctx, command, organizer.ID(), organizer.Name(), organizer.Logo(), organizer.Level()); err != nil {
		log.Println(err)
		return errors.New("failed to create organizer")
	}
	return nil
}

func (s PostgresOrganizerStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.conn.Exec(ctx, "DELETE FROM organizer WHERE id=$1", id)
	return err
}

func (s PostgresOrganizerStorage) Update(ctx context.Context, organizer models.Organizer) error {
	command := "UPDATE organizer SET name = $2, logo = $3, level = $4 WHERE id = $1"

	if _, err := s.conn.Exec(ctx, command, organizer.ID(), organizer.Name(), organizer.Logo(), organizer.Level()); err != nil {
		log.Println(err)
		return errors.New("failed to update")
	}

	return nil
}

func (s PostgresOrganizerStorage) GetLevels(ctx context.Context) ([]models.OrganizerLevel, error) {
	levels := make([]models.OrganizerLevel, 0)

	query := "SELECT * FROM organizer_level"
	rows, err := s.conn.Query(ctx, query)

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
		levels = append(levels, models.NewOrganizerLevel(id, name, code))
	}

	return levels, nil
}

func (s PostgresOrganizerStorage) AddLevel(ctx context.Context, level models.OrganizerLevel) error {
	command := "INSERT INTO organizer_level(id, name, code) VALUES ($1, $2, $3)"

	if _, err := s.conn.Exec(ctx, command, level.ID(), level.Name(), level.Code()); err != nil {
		log.Println(err)
		return errors.New("failed to add to database")
	}
	return nil
}

func (s PostgresOrganizerStorage) DeleteLevel(ctx context.Context, id uuid.UUID) error {
	command := "DELETE FROM organizer_level WHERE id=$1"

	if _, err := s.conn.Exec(ctx, command, id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}
	return nil
}
