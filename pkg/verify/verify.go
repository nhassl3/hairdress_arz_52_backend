package verify

import (
	"github.com/nhassl3/hairdress_arz/pkg/verify/mailer"
)

type Sender struct {
	helper      *Helper
	emailNotify mailer.Notifier
}

func NewSender(codeLen int32, secretKey string, emailNotify mailer.Notifier) *Sender {
	return &Sender{
		helper: NewHelper(
			secretKey,
			codeLen,
		),
		emailNotify: emailNotify,
	}
}

// SendPhone create post request to https://sms.ru/sms/send?api_id=<api_id>&to=<phone>,<phone>&msg=<msg>json=1
func (s *Sender) SendPhone(phone, code string) error {
	return nil
}

// SendEmail compare and send email message with Yandex or Google SMTP server
func (s *Sender) SendEmail(email, code string) error {
	return s.emailNotify.NotifyEmailConfirmation(code, email)
}

func (s *Sender) Helper() *Helper {
	return s.helper
}
