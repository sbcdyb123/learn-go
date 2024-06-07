package repository

import (
	"context"
	"github.com/sbcdyb123/learn-go/internal/repository/cache"
)

var (
	ErrCacheTooFrequently = cache.ErrSetCodeTooFrequently // 验证码发送太频繁
	ErrVerifyTooManyTimes = cache.ErrVerifyTooManyTimes
	ErrUnKnownForCode     = cache.ErrUnKnownForCode
)

type CodeRepository struct {
	Cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		Cache: c,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.Cache.Set(ctx, biz, phone, code)
}
func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.Cache.Verify(ctx, biz, phone, code)
}
