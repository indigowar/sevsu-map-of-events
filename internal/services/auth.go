package services

import (
	"context"

	"github.com/indigowar/map-of-events/internal/domain/services"
)

type authService struct{}

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

func NewAuthService() services.AuthService {
	return &authService{}
}
