package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/pkg/postgres"
)

type PostgresImageStorage struct {
	pool *pgxpool.Pool
}

func (s PostgresImageStorage) GetAllLinks(ctx context.Context) ([]string, error) {
	links := make([]string, 0)

	rows, err := s.pool.Query(ctx, "SELECT link FROM images")
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to read database")
	}

	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			log.Println(err)
			continue
		}
		link := value[0].(string)
		links = append(links, link)
	}

	return links, nil
}

func (s PostgresImageStorage) Get(ctx context.Context, link string) (models.StoredImage, error) {
	query := fmt.Sprintf("SELECT * FROM images WHERE link = '%s'", link)
	var image models.StoredImage
	err := s.pool.QueryRow(ctx, query).Scan(&image.Link, &image.Value)
	if err != nil {
		log.Println(err)
		return models.StoredImage{}, errors.New("failed to find image")
	}
	return image, nil
}

func (s PostgresImageStorage) Add(ctx context.Context, image models.StoredImage) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "INSERT INTO images(link, value) VALUES ($1, $2)"

	if _, err := dataSource.Exec(ctx, command, image.Link, image.Value); err != nil {
		log.Println(err)
		return errors.New("failed to creat image")
	}

	return nil
}

func (s PostgresImageStorage) Remove(ctx context.Context, link string) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "DELETE FROM images WHERE link = $1"

	if _, err := dataSource.Exec(ctx, command, link); err != nil {
		log.Println(err)
		return errors.New("failed to delete image")
	}
	return nil
}

func (s PostgresImageStorage) Update(ctx context.Context, image models.StoredImage) error {
	dataSource := postgres.GetConnectionFromContextOrDefault(ctx, s.pool)

	command := "UPDATE images SET value = $1 WHERE link = $2"

	if _, err := dataSource.Exec(ctx, command, image.Value, image.Link); err != nil {
		log.Println(err)
		return errors.New("failed to update image")
	}
	return nil
}

func NewPostgresImageStorage(pool *pgxpool.Pool) adapters.ImageStorage {
	return &PostgresImageStorage{
		pool: pool,
	}
}
