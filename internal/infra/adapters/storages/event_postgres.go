package storages

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
)

type postgresEventStorage struct {
	con      *pgx.Conn
	subjects storages.SubjectStorageRepository
}

func (s postgresEventStorage) GetIDList(ctx context.Context) ([]uuid.UUID, error) {
	query := "SELECT event_id FROM event"

	results := make([]uuid.UUID, 0)

	rows, err := s.con.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read database")
	}

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to read fetched data from database")
		}
		stringId := vals[0].([16]byte)

		id, err := uuid.FromBytes(stringId[:])
		if err != nil {
			log.Println(err)
			return nil, errors.New("failed to parse id from database")
		}
		results = append(results, id)
	}

	return results, nil
}

func (s postgresEventStorage) Filter(ctx context.Context) ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (s postgresEventStorage) GetByID(ctx context.Context, id uuid.UUID) (models.Event, error) {
	query := fmt.Sprintf("SELECT * FROM event WHERE event_id = '%s'", id.String())

	var Id, organizer, foundingRange, coFoundingRange uuid.UUID
	var title, foundingType, considerationPeriod, realisationPeriod, result, site, document, internalContacts string
	var submissionPeriod time.Time
	var trl int

	err := s.con.QueryRow(ctx, query).Scan(
		&Id, &title, &organizer, &foundingType, &foundingRange, &coFoundingRange, &submissionPeriod,
		&considerationPeriod, &realisationPeriod, &result, &site, &document, &internalContacts,
		&trl,
	)

	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read data from database")
	}

	subjects, err := s.subjects.GetByEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	competitors, err := s.GetCompetitors(ctx, id)
	if err != nil {
		return nil, err
	}

	return models.NewEvent(
		Id,
		title,
		organizer,
		foundingType,
		foundingRange,
		coFoundingRange,
		submissionPeriod,
		considerationPeriod,
		realisationPeriod,
		result,
		site,
		document,
		internalContacts,
		trl,
		competitors,
		subjects,
	), nil
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

	transaction, err := s.con.Begin(ctx)

	if err != nil {
		log.Println(err)
		return errors.New("failed to start a transaction")
	}
	defer transaction.Rollback(ctx)

	_, err = transaction.Exec(ctx, command,
		event.ID(), event.Title(), event.Organizer(), event.FoundingType(), event.FoundingRange(),
		event.CoFoundingRange(), event.SubmissionDeadline(), event.ConsiderationPeriod(), event.RealisationPeriod(),
		event.Result(), event.Site(), event.Document(), event.InternalContacts(), event.TRL())

	if err != nil {
		log.Println(err)
		return errors.New("failed to write in database")
	}
	return nil
}

func (s postgresEventStorage) Delete(ctx context.Context, id uuid.UUID) error {
	command := "DELETE FROM event WHERE event_id = $1"
	if _, err := s.con.Exec(ctx, command, id); err != nil {
		log.Println(err)
		return errors.New("failed to delete value")
	}
	return nil
}

func (s postgresEventStorage) Update(ctx context.Context, event models.Event) error {
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

	_, err := s.con.Exec(ctx, command,
		event.ID(), event.Title(), event.Organizer(), event.FoundingType(), event.FoundingRange(),
		event.CoFoundingRange(), event.SubmissionDeadline(), event.ConsiderationPeriod(), event.RealisationPeriod(),
		event.Result(), event.Site(), event.Document(), event.InternalContacts(), event.TRL())

	if err != nil {
		log.Println(err)
		return errors.New("failed to write in database")
	}
	return nil
}

func (s postgresEventStorage) AddCompetitor(ctx context.Context, id, competitorId uuid.UUID) error {
	command := "INSERT INTO competitor_requirements(cr_id, cr_event, cr_competitor) VALUES ($1, $2, $3)"
	_, err := s.con.Exec(ctx, command, uuid.New(), id, competitorId)
	if err != nil {
		log.Println(err)
		return errors.New("failed to write into database")
	}
	return nil
}

func (s postgresEventStorage) RemoveCompetitor(ctx context.Context, id, competitorId uuid.UUID) error {
	command := "DELETE FROM competitor_requirements WHERE cr_event=$1 AND cr_competitor=$2"
	_, err := s.con.Exec(ctx, command, id, competitorId)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}
	return nil
}

func (s postgresEventStorage) GetCompetitors(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	query := fmt.Sprintf("SELECT cr_competitor FROM competitor_requirements WHERE cr_event='%s'", id.String())

	rows, err := s.con.Query(ctx, query)
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

func NewPostgresEventStorage(con *pgx.Conn, subjects storages.SubjectStorageRepository) storages.EventStorageRepository {
	return &postgresEventStorage{
		con:      con,
		subjects: subjects,
	}
}
