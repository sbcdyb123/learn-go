package retryable

import (
	"context"
	"github.com/sbcdyb123/learn-go/internal/service/sms"
)

type Service struct {
	svc      sms.Service
	retryCnt int
}

func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tplId, args, numbers...)
	for err != nil && s.retryCnt < 10 {
		err = s.svc.Send(ctx, tplId, args, numbers...)
		s.retryCnt++
		return err
	}
	return err
}
