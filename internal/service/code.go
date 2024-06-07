package service

import (
	"context"
	"fmt"
	"github.com/sbcdyb123/learn-go/internal/repository"
	"github.com/sbcdyb123/learn-go/internal/service/sms"
	"math/rand"
)

const codeTplId = "1877556"

var (
	ErrCacheTooFrequently = repository.ErrCacheTooFrequently
)

type CodeService struct {
	repo   *repository.CodeRepository
	SmsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		SmsSvc: smsSvc,
	}
}

func (s *CodeService) Send(ctx context.Context, biz string, phone string) error {
	// 发送验证码
	code := s.generateCode()
	fmt.Println("验证码:", code)
	err := s.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = s.SmsSvc.Send(ctx, codeTplId, []string{}, phone)
	return err
}
func (s *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return s.repo.Verify(ctx, biz, phone, inputCode)
}

func (s *CodeService) generateCode() string {
	// 生成随机验证码
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
