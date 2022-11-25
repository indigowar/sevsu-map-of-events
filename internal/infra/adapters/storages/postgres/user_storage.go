package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
)

type userStorage struct {
	pool *pgxpool.Pool
}

func (storage userStorage) GetByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (storage userStorage) GetByName(ctx context.Context, name string) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (storage userStorage) Create(ctx context.Context, user models.User) error {
	//TODO implement me
	panic("implement me")
}

func (storage userStorage) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (storage userStorage) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
	//TODO implement me
	panic("implement me")
}

func (storage userStorage) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresUserStorage(pool *pgxpool.Pool) storages.UserStorage {
	return &userStorage{
		pool: pool,
	}
}
