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

type coFoundingRangePostgresStorage struct {
	con *pgx.Conn
}

func (s coFoundingRangePostgresStorage) GetByID(ctx context.Context, id uuid.UUID) (models.RangeModel, error) {
	var low, high int

	query := fmt.Sprintf("SELECT co_founding_low, co_founding_high FROM co_founding_range WHERE co_founding_range_id = '%s'", id.String())

	if err := s.con.QueryRow(ctx, query).Scan(&low, &high); err != nil {
		log.Println("Got query error or scan error: ", err)
		return models.RangeModel{}, err
	}

	return models.RangeModel{ID: id, Low: low, High: high}, nil
}

func (s coFoundingRangePostgresStorage) GetMaximumRange(ctx context.Context) (models.RangeModel, error) {
	var low, high int
	if s.con.QueryRow(ctx, "SELECT MIN(co_founding_low) FROM co_founding_range").Scan(&low) != nil ||
		s.con.QueryRow(ctx, "SELECT MAX(co_founding_high) FROM co_founding_range").Scan(&high) != nil {
		return models.RangeModel{}, errors.New("failed to read database")
	}
	return models.RangeModel{Low: low, High: high}, nil
}

func (s coFoundingRangePostgresStorage) Create(ctx context.Context, foundingRange models.RangeModel) (models.RangeModel, error) {
	command := "INSERT INTO co_founding_range (co_founding_range_id, co_founding_low, co_founding_high) VALUES ($1, $2, $3)"
	_, err := s.con.Exec(ctx, command, foundingRange.ID, foundingRange.Low, foundingRange.High)
	if err != nil {
		log.Println(err)
		return models.RangeModel{}, errors.New("failed to insert")
	}

	return s.GetByID(ctx, foundingRange.ID)
}

func (s coFoundingRangePostgresStorage) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.con.Exec(ctx, "DELETE FROM co_founding_range WHERE co_founding_range_id = $1", id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}

	return nil
}

func NewCoFoundingRangePostgresStorage(con *pgx.Conn) storages.RangeStorageRepository {
	return &coFoundingRangePostgresStorage{
		con: con,
	}
}
