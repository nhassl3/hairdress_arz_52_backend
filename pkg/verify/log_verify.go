package sms

import (
	"context"

	"github.com/nhassl3/hairdress_arz/pkg/verify"
	"go.uber.org/zap"
)

type SenderLog struct {
	log    *zap.Logger
	helper *verify.Helper
}

func NewSenderLog(logger *zap.Logger, codeWidth int32, secretKey string) *SenderLog {
	return &SenderLog{
		log: logger,
		helper: verify.NewHelper(
			secretKey,
			codeWidth,
		),
	}
}

func (s *SenderLog) SendPhone(ctx context.Context, phone, code string) error {
	_ = ctx
	s.log.Info("sms code", zap.String("phone", phone), zap.String("code", code))
	return nil
}

func (s *SenderLog) SendEmail(ctx context.Context, email, code string) error {
	_ = ctx
	s.log.Info("email code", zap.String("email", email), zap.String("code", code))
	return nil
}

func (s *SenderLog) Helper() *verify.Helper {
	return s.helper
}
