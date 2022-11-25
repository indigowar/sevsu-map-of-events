package services

import (
	"context"

	"github.com/indigowar/map-of-events/internal/config"
	"github.com/indigowar/map-of-events/internal/domain/adapters"
	"github.com/indigowar/map-of-events/internal/domain/services"
)

type authService struct {
	userStorage adapters.UserStorage
	config      config.AuthConfig
}

func (svc authService) Login(ctx context.Context, name, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (svc authService) GetRefresh(ctx context.Context, rt string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (svc authService) GetAccess(ctx context.Context, rt string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (svc authService) CreateUser(ctx context.Context, name, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthService(userStorage adapters.UserStorage, cfg config.AuthConfig) services.AuthService {
	return &authService{
		userStorage: userStorage,
		config:      cfg,
	}
}
