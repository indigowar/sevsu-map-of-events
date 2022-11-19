package services

import (
	"context"
	"errors"
	"log"

	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/repos/adapters/storages"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type imageService struct {
	storage storages.ImageStorage
}

func (svc imageService) Get(ctx context.Context, link string) (models.StoredImage, error) {
	image, err := svc.storage.Get(ctx, link)
	if err != nil {
		log.Println(err)
		return models.StoredImage{}, errors.New("failed to read database")
	}
	return image, nil
}

func (svc imageService) Create(ctx context.Context, link string, image []byte) (models.StoredImage, error) {
	model := models.StoredImage{
		Link:  link,
		Value: image,
	}

	err := svc.storage.Create(ctx, model)
	if err != nil {
		log.Println(err)
		return models.StoredImage{}, errors.New("failed to add image")
	}

	return svc.Get(ctx, link)
}

func (svc imageService) Delete(ctx context.Context, link string) error {
	err := svc.storage.Delete(ctx, link)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete image")
	}
	return nil
}

func (svc imageService) Update(ctx context.Context, link string, image []byte) (models.StoredImage, error) {
	model := models.StoredImage{
		Link:  link,
		Value: image,
	}

	err := svc.storage.Update(ctx, model)
	if err != nil {
		log.Println(err)
		return models.StoredImage{}, errors.New("failed to update an image")
	}

	return svc.Get(ctx, link)
}

func NewImageService(storage storages.ImageStorage) services.ImageService {
	return &imageService{
		storage: storage,
	}

}
