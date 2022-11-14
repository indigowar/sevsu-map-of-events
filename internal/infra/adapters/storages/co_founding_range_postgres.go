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

func (s coFoundingRangePostgresStorage) GetByID(ctx context.Context, id uuid.UUID) (models.FoundingRange, error) {
	var Id uuid.UUID
	var low, high int

	query := fmt.Sprintf("SELECT * FROM co_founding_range WHERE co_founding_range_id = '%s'", id.String())

	if err := s.con.QueryRow(ctx, query).Scan(&Id, &low, &high); err != nil {
		log.Println("Got query error or scan error: ", err)
		return nil, err
	}

	return models.NewRange(Id, low, high), nil
}

func (s coFoundingRangePostgresStorage) GetMaximumRange(ctx context.Context) (models.FoundingRange, error) {
	var low, high int
	if s.con.QueryRow(ctx, "SELECT MIN(co_founding_low) FROM co_founding_range").Scan(&low) != nil ||
		s.con.QueryRow(ctx, "SELECT MAX(co_founding_high) FROM co_founding_range").Scan(&high) != nil {
		return nil, errors.New("failed to read database")
	}
	return models.NewRange(uuid.UUID{}, low, high), nil
}

func (s coFoundingRangePostgresStorage) Create(ctx context.Context, foundingRange models.FoundingRange) (models.FoundingRange, error) {
	command := "INSERT INTO co_founding_range (co_founding_range_id, co_founding_low, co_founding_high) VALUES ($1, $2, $3)"
	_, err := s.con.Exec(ctx, command, foundingRange.ID(), foundingRange.Low(), foundingRange.High())
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to insert")
	}

	return s.GetByID(ctx, foundingRange.ID())
}

func (s coFoundingRangePostgresStorage) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.con.Exec(ctx, "DELETE FROM co_founding_range WHERE co_founding_id = $1", id); err != nil {
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
