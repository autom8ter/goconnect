package grid

import (
	"github.com/autom8ter/goconnect/pkg/errors"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/grpc/codes"
)

type SendGrid struct {
	client *sendgrid.Client
}

func NewSendGrid(client *sendgrid.Client) *SendGrid {
	return &SendGrid{
		client: client,
	}
}

func (s *SendGrid) Send(from *contacts.Contact, to *contacts.Contact, subject, plaintext, htmlAlternative string) error {
	fromemail := mail.NewEmail(from.Name, from.Email)
	toemail := mail.NewEmail(to.Name, to.Email)
	message := mail.NewSingleEmail(fromemail, subject, toemail, plaintext, htmlAlternative)
	_, err := s.client.Send(message)
	if err != nil {
		return errors.StatusStack(err, codes.Internal, "server.UpdateAccount: %s", "failed to update account")
	}
	return nil
}
