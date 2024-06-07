package repository

import (
	"context"
	"database/sql"
	"github.com/sbcdyb123/learn-go/internal/domain"
	"github.com/sbcdyb123/learn-go/internal/repository/cache"
	"github.com/sbcdyb123/learn-go/internal/repository/dao"
	"log"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
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
	return r.dao.Inset(ctx, r.entityToDao(u))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
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
	u = r.entityToDomain(u1)
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			log.Println("cache set error:", err)
		}
	}()
	return u, err
}
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Username: u.Username,
		BirthDay: u.BirthDay,
		Intro:    u.Intro,
		CTime:    time.UnixMilli(u.UTime),
		UTime:    time.UnixMilli(u.CTime),
	}
}
func (r *UserRepository) entityToDao(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Username: u.Username,
		BirthDay: u.BirthDay,
		Intro:    u.Intro,
		CTime:    u.CTime.UnixMilli(),
		UTime:    u.UTime.UnixMilli(),
	}
}
