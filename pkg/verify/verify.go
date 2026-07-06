package sms

import (
	"context"

	"github.com/nhassl3/hairdress_arz/pkg/verify"
)

type Sender struct {
	helper *verify.Helper
}

func NewSender(codeWidth int32, secretKey string) *Sender {
	return &Sender{
		helper: verify.NewHelper(
			secretKey,
			codeWidth,
		),
	}
}

// SendPhone create post request to https://sms.ru/sms/send?api_id=<api_id>&to=<phone>,<phone>&msg=<msg>json=1
func (s *Sender) SendPhone(ctx context.Context, phone, code string) error {
	return nil
}

// SendEmail compare email message
func (s *Sender) SendEmail(ctx context.Context, email, code string) error {
	return nil
}

func (s *Sender) Helper() *verify.Helper {
	return s.helper
}
