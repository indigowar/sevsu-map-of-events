package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/pkg/errors"
	"github.com/indigowar/map-of-events/pkg/postgres"
)

type competitorStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresCompetitorStorage(p *pgxpool.Pool) adapters.CompetitorStorage {
	return &competitorStorage{pool: p}
}

func (s competitorStorage) AllIDs(ctx context.Context) ([]uuid.UUID, errors.Error) {
	query := "SELECT competitor_id FROM competitor"

	result := make([]uuid.UUID, 0)

	rows, err := s.pool.Query(ctx, query)

	if err != nil {
		return nil, handleTheError(err, createInternalStorageError(err, "failed to read database"))
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, createInternalStorageError(err, "failed to read values of a row")
		}

		id, e := parseIDFromValue(values[0])
		if e != nil {
			return nil, e
		}
		result = append(result, id)
	}

	return result, nil
}

func (s competitorStorage) Get(ctx context.Context, id uuid.UUID) (models.Competitor, errors.Error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)
	var c models.Competitor
	query := fmt.Sprintf("SELECT * FROM competitor WHERE competitor_id == '%s'", id.String())
	if err := dataSource.QueryRow(ctx, query).Scan(&c.ID, &c.Name); err != nil {
		return models.Competitor{}, handleTheError(err,
			errors.CreateError(adapters.ErrReasonObjectNotFoundErr, "object was not found",
				fmt.Sprintf("object %s was not found in competitors", id.String())))
	}
	return c, nil
}

func (s competitorStorage) GetAll(ctx context.Context) ([]models.Competitor, errors.Error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	comps := make([]models.Competitor, 0)

	rows, err := dataSource.Query(ctx, "SELECT * FROM competitor")
	if err != nil {
		return nil, handleTheError(err, createInternalStorageError(err, "failed to read database"))
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, handleTheError(err, createInternalStorageError(err, "failed to read values"))
		}
		id, e := parseIDFromValue(values[0])
		if e != nil {
			return nil, e
		}
		name := values[1].(string)
		comps = append(comps, models.Competitor{ID: id, Name: name})
	}

	return comps, nil
}

func (s competitorStorage) Create(ctx context.Context, competitor models.Competitor) errors.Error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)
	command := "INSERT INTO competitor (competitor_id, competitor_name) VALUES ($1, $2)"
	if _, err := dataSource.Exec(ctx, command, competitor.ID, competitor.Name); err != nil {
		return createInternalStorageError(err, "failed to create a competitor")
	}
	return nil
}

func (s competitorStorage) Update(ctx context.Context, competitor models.Competitor) errors.Error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)
	command := "UPDATE competitor SET competitor_name = $2 WHERE competitor_id = $1"
	if _, err := dataSource.Exec(ctx, command, competitor.ID, competitor.Name); err != nil {
		return createInternalStorageError(err, "failed to update a competitor")
	}
	return nil
}

func (s competitorStorage) Delete(ctx context.Context, id uuid.UUID) errors.Error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)
	_, err := dataSource.Exec(ctx, "DELETE FROM competitor WHERE competitor_id = $1", id)
	return createInternalStorageError(err, "failed to delete a competitor")
}

func (s competitorStorage) BeginTransaction(ctx context.Context) (interface{}, error) {
	return s.pool.Begin(ctx)
}

func (s competitorStorage) CloseTransaction(ctx context.Context, transaction interface{}) error {
	tx := transaction.(*pgxpool.Tx)
	return tx.Rollback(ctx)
}

// just a shortcut of errors.CreateError(storages.ErrReasonInternalStorageErr, ...)
func createInternalStorageError(e error, failedJob string) errors.Error {
	return errors.CreateError(adapters.ErrReasonInternalStorageErr, failedJob, failedJob+":"+e.Error())
}

// will check all errors, if they're about transaction will return the transaction failure error
func handleTheError(e error, onRows errors.Error) errors.Error {
	switch e {
	case pgx.ErrNoRows:
		return onRows
	case pgx.ErrTxClosed:
	case pgx.ErrTxCommitRollback:
		return createInternalStorageError(e, "transaction failure")
	}
	return nil
}

// parse from interface{} ID
// under interface{} should be [16]byte
func parseIDFromValue(v interface{}) (uuid.UUID, errors.Error) {
	byteId := v.([16]byte)
	id, err := uuid.FromBytes(byteId[:])
	if err != nil {
		return uuid.UUID{}, createInternalStorageError(err, "failed to parse id from a row")
	}
	return id, nil
}
