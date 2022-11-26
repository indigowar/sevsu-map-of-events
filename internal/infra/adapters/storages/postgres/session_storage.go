package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/models"
)

type sessionStorage struct {
	pool *pgxpool.Pool
}

func (s sessionStorage) GetByToken(ctx context.Context, token string) (models.TokenSession, error) {
	//TODO implement me
	panic("implement me")
}

func (s sessionStorage) GetByUser(ctx context.Context, id uuid.UUID) (models.TokenSession, error) {
	//TODO implement me
	panic("implement me")
}

func (s sessionStorage) Create(ctx context.Context, session models.TokenSession) error {
	//TODO implement me
	panic("implement me")
}

func (s sessionStorage) Delete(ctx context.Context, token string) error {
	//TODO implement me
	panic("implement me")
}

func NewSessionStorage(p *pgxpool.Pool) adapters.SessionStorage {
	return &sessionStorage{
		pool: p,
	}
}
