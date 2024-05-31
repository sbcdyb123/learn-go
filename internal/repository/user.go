package repository

import (
	"context"
	"github.com/sbcdyb123/learn-go/internal/domain"
	"github.com/sbcdyb123/learn-go/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Inset(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
func (r *UserRepository) Updates(ctx context.Context, u domain.User) error {
	err := r.dao.Updates(ctx, dao.User{
		Id:       u.Id,
		Username: u.Username,
		BirthDay: u.BirthDay,
		Intro:    u.Intro,
	})
	return err
}
