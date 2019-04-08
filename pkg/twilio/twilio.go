package twilio

import (
	"github.com/autom8ter/goconnect/pkg/errors"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"github.com/sfreiberg/gotwilio"
	"google.golang.org/grpc/codes"
)

type Twilio struct {
	client *gotwilio.Twilio
}

func New(client *gotwilio.Twilio) *Twilio {
	return &Twilio{client}
}

func (t *Twilio) SendSMS(from, to *contacts.Contact, body, callback, appsid string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := t.client.SendSMS(from.Phone, to.Phone, body, callback, appsid)
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to send sms")
	}
	return resp, nil
}

func (t *Twilio) SendMMS(from, to *contacts.Contact, body, mediaurl, callback, appsid string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := t.client.SendMMS(from.Phone, to.Phone, body, mediaurl, callback, appsid)
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to send mms")
	}
	return resp, nil
}

func (t *Twilio) Call(from, to *contacts.Contact, callback string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := t.client.CallWithUrlCallbacks(from.Phone, to.Phone, gotwilio.NewCallbackParameters(callback))
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to send call")
	}
	return resp, nil
}
