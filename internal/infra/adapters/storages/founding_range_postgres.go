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

type foundingRangeStorage struct {
	con *pgx.Conn
}

func (s foundingRangeStorage) GetByID(ctx context.Context, id uuid.UUID) (models.FoundingRange, error) {
	var Id uuid.UUID
	var low, high int

	query := fmt.Sprintf("SELECT * FROM founding_range WHERE id = '%s'", id.String())

	if err := s.con.QueryRow(ctx, query).Scan(&Id, &low, &high); err != nil {
		log.Println("Got query error or scan error: ", err)
		return nil, err
	}

	return models.NewRange(Id, low, high), nil
}

func (s foundingRangeStorage) GetMaximumRange(ctx context.Context) (models.FoundingRange, error) {
	var low, high int
	if s.con.QueryRow(ctx, "SELECT MIN(low) FROM founding_range").Scan(&low) != nil ||
		s.con.QueryRow(ctx, "SELECT MAX(high) FROM founding_range").Scan(&high) != nil {
		return nil, errors.New("failed to read database")
	}
	return models.NewRange(uuid.UUID{}, low, high), nil
}

func (s foundingRangeStorage) Create(ctx context.Context, foundingRange models.FoundingRange) (models.FoundingRange, error) {
	command := "INSERT INTO founding_range (id, low, high) VALUES ($1, $2, $3)"

	_, err := s.con.Exec(ctx, command, foundingRange.ID(), foundingRange.Low(), foundingRange.High())
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to insert")
	}

	return s.GetByID(ctx, foundingRange.ID())
}

func (s foundingRangeStorage) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.con.Exec(ctx, "DELETE FROM founding_range WHERE id = $1", id); err != nil {
		log.Println(err)
		return errors.New("failed to delete from database")
	}

	return nil
}

func NewFoundingRangePostgresStorage(con *pgx.Conn) storages.RangeStorageRepository {
	return &foundingRangeStorage{
		con: con,
	}
}
