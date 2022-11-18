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

type foundingRangeStorage struct {
	pool *pgxpool.Pool
}

func (s foundingRangeStorage) InvokeTransactionMechanism(ctx context.Context) (interface{}, error) {
	return s.pool.Begin(ctx)
}

func (s foundingRangeStorage) ShadowTransactionMechanism(ctx context.Context, transaction interface{}) error {
	tx := transaction.(*pgxpool.Tx)
	return tx.Rollback(ctx)
}

func (s foundingRangeStorage) GetByID(ctx context.Context, id uuid.UUID) (models.RangeModel, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	var Id uuid.UUID
	var low, high int

	query := fmt.Sprintf("SELECT * FROM founding_range WHERE founding_range_id = '%s'", id.String())

	if err := dataSource.QueryRow(ctx, query).Scan(&Id, &low, &high); err != nil {
		log.Println("Got query error or scan error: ", err)
		return models.RangeModel{}, err
	}

	return models.RangeModel{ID: Id, Low: low, High: high}, nil
}

func (s foundingRangeStorage) GetMaximumRange(ctx context.Context) (models.RangeModel, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	var low, high int
	if dataSource.QueryRow(ctx, "SELECT MIN(founding_range_low) FROM founding_range").Scan(&low) != nil ||
		dataSource.QueryRow(ctx, "SELECT MAX(founding_range_high) FROM founding_range").Scan(&high) != nil {
		return models.RangeModel{}, errors.New("failed to read database")
	}
	return models.RangeModel{Low: low, High: high}, nil
}

func (s foundingRangeStorage) Create(ctx context.Context, foundingRange models.RangeModel) (models.RangeModel, error) {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "INSERT INTO founding_range (founding_range_id, founding_range_low, founding_range_high) VALUES ($1, $2, $3)"

	_, err := dataSource.Exec(ctx, command, foundingRange.ID, foundingRange.Low, foundingRange.High)
	if err != nil {
		log.Println(err)
		return models.RangeModel{}, errors.New("failed to insert")
	}

	return s.GetByID(ctx, foundingRange.ID)
}

func (s foundingRangeStorage) Delete(ctx context.Context, id uuid.UUID) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	if _, err := dataSource.Exec(ctx, "DELETE FROM founding_range WHERE founding_range_id = $1", id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}

	return nil
}

func NewFoundingRangePostgresStorage(p *pgxpool.Pool) storages.RangeStorageRepository {
	return &foundingRangeStorage{
		pool: p,
	}
}
