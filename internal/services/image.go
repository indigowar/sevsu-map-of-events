package services

import (
	"context"
	"errors"
	"log"

	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/models"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type imageService struct {
	storage adapters.ImageStorage
}

func (svc imageService) GetAllLinks(ctx context.Context) ([]string, error) {
	links, err := svc.storage.GetAllLinks(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return links, nil
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

	err := svc.storage.Add(ctx, model)
	if err != nil {
		log.Println(err)
		return models.StoredImage{}, errors.New("failed to add image")
	}

	return model, nil
}

func (svc imageService) Delete(ctx context.Context, link string) error {
	err := svc.storage.Remove(ctx, link)
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

func NewImageService(storage adapters.ImageStorage) services.ImageService {
	return &imageService{
		storage: storage,
	}

}
