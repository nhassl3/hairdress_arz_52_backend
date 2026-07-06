package verify

import (
	"go.uber.org/zap"
)

type SenderLog struct {
	log    *zap.Logger
	helper *Helper
}

func NewSenderLog(logger *zap.Logger, codeWidth int32, secretKey string) *SenderLog {
	return &SenderLog{
		log: logger,
		helper: NewHelper(
			secretKey,
			codeWidth,
		),
	}
}

func (s *SenderLog) SendPhone(phone, code string) error {
	s.log.Info("sms code", zap.String("phone", phone), zap.String("code", code))
	return nil
}

func (s *SenderLog) SendEmail(email, code string) error {
	s.log.Info("email code", zap.String("email", email), zap.String("code", code))
	return nil
}

func (s *SenderLog) Helper() *Helper {
	return s.helper
}
