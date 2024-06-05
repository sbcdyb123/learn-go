package repository

import (
	"context"
	"github.com/sbcdyb123/learn-go/internal/domain"
	"github.com/sbcdyb123/learn-go/internal/repository/cache"
	"github.com/sbcdyb123/learn-go/internal/repository/dao"
	"log"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
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
func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从缓存中查找
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	u1, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	// 更新缓存
	u = domain.User{
		Id:       u1.Id,
		Email:    u1.Email,
		Password: u1.Password,
		Username: u1.Username,
		BirthDay: u1.BirthDay,
		Intro:    u1.Intro,
	}
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			log.Println("cache set error:", err)
		}
	}()
	return u, err
}
