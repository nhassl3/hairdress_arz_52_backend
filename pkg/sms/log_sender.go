package sms

import (
	"context"

	"go.uber.org/zap"
)

type SMSenderLog struct {
	log    *zap.Logger
	helper *Helper
}

func NewSMSSenderLog(logger *zap.Logger, codeWidth int32, secretKey string) *SMSenderLog {
	return &SMSenderLog{
		log: logger,
		helper: NewHelper(
			secretKey,
			codeWidth,
		),
	}
}

func (s *SMSenderLog) Send(ctx context.Context, phone, code string) error {
	_ = ctx
	s.log.Info("sms code", zap.String("phone", phone), zap.String("code", code))
	return nil
}

func (s *SMSenderLog) Helper() *Helper {
	return s.helper
}
