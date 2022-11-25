package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/models"
)

type userStorage struct {
	pool *pgxpool.Pool
}

func (storage userStorage) GetByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM authenticatedusers WHERE id = '%s'", id.String())

	var user models.User

	err := storage.pool.QueryRow(ctx, query).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		log.Println(err)
		return models.User{}, errors.New("user was not found")
	}

	return user, nil
}

func (storage userStorage) GetByName(ctx context.Context, name string) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM authenticatedusers WHERE name= '%s'", name)

	var user models.User

	err := storage.pool.QueryRow(ctx, query).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		log.Println(err)
		return models.User{}, errors.New("user was not found")
	}
	return user, nil
}

func (storage userStorage) Create(ctx context.Context, user models.User) error {
	command := "INSERT INTO authenticatedusers(id, name, password) VALUES ($1, $2, $3)"

	_, err := storage.pool.Exec(ctx, command, user.ID, user.Name, user.Password)
	if err != nil {
		log.Println(err)
		return errors.New("failed to add user")
	}
	return nil
}

func (storage userStorage) Delete(ctx context.Context, id uuid.UUID) error {
	command := "DELETE FROM authenticatedusers WHERE id = $1"

	_, err := storage.pool.Exec(ctx, command, id.String())
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete a user")
	}
	return nil
}

func (storage userStorage) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
	command := "UPDATE authenticatedUsers SET name = $2 WHERE id = $1"

	_, err := storage.pool.Exec(ctx, command, id.String(), name)
	if err != nil {
		log.Println(err)
		return errors.New("failed to update")
	}
	return nil
}

func (storage userStorage) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	command := "UPDATE authenticatedUsers SET password = $1 WHERE id = $2"
	_, err := storage.pool.Exec(ctx, command, password, id)
	if err != nil {
		log.Println(err)
		return errors.New("failed to update user")
	}
	return nil
}

func NewPostgresUserStorage(pool *pgxpool.Pool) adapters.UserStorage {
	return &userStorage{
		pool: pool,
	}
}
