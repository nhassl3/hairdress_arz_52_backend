package sms

import "context"

type SMSender struct {
	helper *Helper
}

func NewSMSender(codeWidth int32, secretKey string) *SMSender {
	return &SMSender{
		helper: NewHelper(
			secretKey,
			codeWidth,
		),
	}
}

// Send create post request to https://sms.ru/sms/send?api_id=<api_id>&to=<phone>,<phone>&msg=<msg>json=1
func (s *SMSender) Send(ctx context.Context, phone, code string) error {
	return nil
}

func (s *SMSender) Helper() *Helper {
	return s.helper
}
