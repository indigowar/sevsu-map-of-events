package storages

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/pkg/postgres"
)

type postgresEventStorage struct {
	pool *pgxpool.Pool
}

func (s postgresEventStorage) InvokeTransactionMechanism(ctx context.Context) (interface{}, error) {
	return s.pool.Begin(ctx), nil
}

func (s postgresEventStorage) ShadowTransactionMechanism(ctx context.Context, transaction interface{}) error {
	tx := transaction.(*pgxpool.Tx)
	return tx.Rollback(ctx)
}

func (s postgresEventStorage) GetIDList(ctx context.Context) ([]uuid.UUID, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	query := "SELECT event_id FROM event"

	results := make([]uuid.UUID, 0)

	rows, err := dataSource.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read database")
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read fetched data from database")
		}
		stringId := values[0].([16]byte)

		id, err := uuid.FromBytes(stringId[:])
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to parse id from database")
		}
		results = append(results, id)
	}

	return results, nil
}

func (s postgresEventStorage) Filter(_ context.Context) ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Event, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	query := fmt.Sprintf("SELECT * FROM event WHERE event_id = '%s'", id.String())

	var event models.Event

	err := dataSource.QueryRow(ctx, query).Scan(
		&event.ID, &event.Title, &event.Organizer, &event.FoundingType, &event.FoundingRange, &event.CoFoundingRange, &event.SubmissionDeadline,
		&event.ConsiderationPeriod, &event.RealisationPeriod, &event.Result, &event.Site, &event.Document, &event.InternalContacts,
		&event.TRL,
	)

	if err != nil {
		log.Println(err)
		return models.Event{}, errors.New("failed to read data from database")
	}

	event.Competitors, err = s.GetCompetitors(ctx, id)
	if err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func (s postgresEventStorage) Create(ctx context.Context, event models.Event) error {
	command :=
		`INSERT INTO event(
                  event_id, title, event_organizer, event_founding_type, event_founding_range, event_co_founding_range,
                  event_submission_deadline, event_consideration_period, event_realisation_period, event_result,
                  event_site, event_document, event_internal_contacts, event_trl)
		VALUES 
		(
		 $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`

	transaction, err := s.pool.Begin(ctx)

	if err != nil {
		log.Println(err)
		return errors.New("failed to start a transaction")
	}
	defer func(transaction pgx.Tx, ctx context.Context) {
		_ = transaction.Rollback(ctx)
	}(transaction, ctx)

	_, err = transaction.Exec(ctx, command,
		event.ID, event.Title, event.Organizer, event.FoundingType, event.FoundingRange,
		event.CoFoundingRange, event.SubmissionDeadline, event.ConsiderationPeriod, event.RealisationPeriod,
		event.Result, event.Site, event.Document, event.InternalContacts, event.TRL)

	if err != nil {
		log.Println(err)
		return errors.New("failed to write in database")
	}
	return nil
}

func (s postgresEventStorage) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("transaction failure")
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM event WHERE event_id = $1", id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from db")
	}

	if _, err := tx.Exec(ctx, "DELETE FROM competitor_requirements WHERE cr_event = $1", id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from db")
	}

	return nil
}

func (s postgresEventStorage) Update(ctx context.Context, event models.Event) error {
	// TODO: Add updating of competitor requirements
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("failed to start a transaction")
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	command := ` UPDATE event SET 
                  title = $1,
                  event_organizer = $2,
                  event_founding_type = $3,
                  event_founding_range = $4,
                  event_co_founding_range = $5,
                  event_submission_deadline = $6,
                  event_consideration_period = $7,
                  event_realisation_period = $8,
                  event_result = $9,
                  event_site = $10,
                  event_document = $11,
                  event_internal_contacts = $12,
                  event_trl = $13
                  WHERE event_id = $1`

	_, err = tx.Exec(ctx, command,
		event.ID, event.Title, event.Organizer, event.FoundingType, event.FoundingRange,
		event.CoFoundingRange, event.SubmissionDeadline, event.ConsiderationPeriod, event.RealisationPeriod,
		event.Result, event.Site, event.Document, event.InternalContacts, event.TRL)

	if err != nil {
		log.Println(err)
		return errors.New("failed to write in database")
	}
	return nil
}

func (s postgresEventStorage) AddCompetitor(ctx context.Context, id, competitorId uuid.UUID) error {
	command := "INSERT INTO competitor_requirements(cr_id, cr_event, cr_competitor) VALUES ($1, $2, $3)"
	_, err := s.pool.Exec(ctx, command, uuid.New(), id, competitorId)
	if err != nil {
		log.Println(err)
		return errors.New("failed to write into database")
	}
	return nil
}

func (s postgresEventStorage) RemoveCompetitor(ctx context.Context, id, competitorId uuid.UUID) error {
	command := "DELETE FROM competitor_requirements WHERE cr_event=$1 AND cr_competitor=$2"
	_, err := s.pool.Exec(ctx, command, id, competitorId)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}
	return nil
}

func (s postgresEventStorage) GetCompetitors(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	query := fmt.Sprintf("SELECT cr_competitor FROM competitor_requirements WHERE cr_event='%s'", id.String())

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read from database")
	}

	ids := make([]uuid.UUID, 0)

	for rows.Next() {
		val, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read fetched data")
		}

		byteId := val[0].([16]byte)

		id, err := uuid.FromBytes(byteId[:])
		if err != nil {
			return nil, errors.New("failed to parse data")
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func NewPostgresEventStorage(p *pgxpool.Pool) storages.EventStorageRepository {
	return &postgresEventStorage{
		pool: p,
	}
}
