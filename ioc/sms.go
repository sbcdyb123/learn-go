package ioc

import (
	"github.com/sbcdyb123/learn-go/internal/service/sms"
	"github.com/sbcdyb123/learn-go/internal/service/sms/memory"
)

func InitSmsService() sms.Service {
	return memory.NewService()
}
