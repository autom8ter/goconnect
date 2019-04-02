package email

import "github.com/sendgrid/sendgrid-go/helpers/mail"

type EmailOption func(m *mail.SGMailV3)

func NewEmail(opts ...EmailOption) *mail.SGMailV3 {
	m := &mail.SGMailV3{}
	for _, o := range opts {
		o(m)
	}
	return m
}
