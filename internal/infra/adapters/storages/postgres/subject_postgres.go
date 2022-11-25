package postgres

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

type postgresSubjectStorage struct {
	pool *pgxpool.Pool
}

func (s postgresSubjectStorage) BeginTransaction(ctx context.Context) (interface{}, error) {
	return s.pool.Begin(ctx)
}

func (s postgresSubjectStorage) CloseTransaction(ctx context.Context, transaction interface{}) error {
	tx := transaction.(*pgxpool.Tx)
	return tx.Rollback(ctx)
}

func (s postgresSubjectStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Subject, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	query := fmt.Sprintf("SELECT * FROM subject WHERE subject_id = '%s'", id.String())

	subject := models.Subject{
		ID: id,
	}

	if err := dataSource.QueryRow(ctx, query).Scan(&subject.ID, &subject.Name, &subject.Name); err != nil {
		log.Println(err)
		return models.Subject{}, errors.New("failed to read from database")
	}

	return subject, nil
}

func (s postgresSubjectStorage) GetByEvent(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	query := fmt.Sprintf("SELECT subject_id FROM subject WHERE subject_event = '%s'", id.String())

	rows, err := dataSource.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read data from database")
	}

	result := make([]uuid.UUID, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read data from database")
		}

		byteId := values[0].([16]byte)
		id, err := uuid.FromBytes(byteId[:])
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read data from database")
		}

		result = append(result, id)
	}

	return result, nil
}

func (s postgresSubjectStorage) GetAll(ctx context.Context) ([]models.Subject, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	rows, err := dataSource.Query(ctx, "SELECT subject_id, subject_name FROM subject")
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read data")
	}

	subjects := make([]models.Subject, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to data")
		}
		parsedId := values[0].([16]byte)
		id, err := uuid.FromBytes(parsedId[:])
		if err != nil {
			log.Println(err)
			continue
		}
		name := values[1].(string)

		subjects = append(subjects, models.Subject{ID: id, Name: name})
	}

	return subjects, nil
}

func (s postgresSubjectStorage) Add(ctx context.Context, subject models.Subject) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "INSERT INTO subject(subject_id, subject_name, subject_event) VALUES ($1, $2, $3)"
	_, err := dataSource.Exec(ctx, command, subject.ID, subject.Name, subject.EventID)
	if err != nil {
		log.Println(err)
		return errors.New("failed to write in the database")
	}
	return nil
}

func (s postgresSubjectStorage) Delete(ctx context.Context, id uuid.UUID) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "DELETE FROM subject WHERE subject_name = $1"

	if _, err := dataSource.Exec(ctx, command, id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}
	return nil
}

func (s postgresSubjectStorage) Update(ctx context.Context, subject models.Subject) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "UPDATE subject SET subject_name = $1, subject_event = $2 WHERE subject_id = $3"

	if _, err := dataSource.Exec(ctx, command, subject.Name, subject.EventID, subject.ID); err != nil {
		log.Println(err)
		return errors.New("failed to update the database")
	}
	return nil
}

func NewPostgresSubjectStorage(p *pgxpool.Pool) storages.SubjectStorage {
	return &postgresSubjectStorage{
		pool: p,
	}
}
