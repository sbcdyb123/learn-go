package service

import (
	"context"
	"github.com/sbcdyb123/learn-go/internal/domain"
	"github.com/sbcdyb123/learn-go/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	return svc.repo.Create(ctx, user)
}
