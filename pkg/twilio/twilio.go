package twilio

import (
	"github.com/autom8ter/goconnect/pkg/errors"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"github.com/autom8ter/objectify"
	"github.com/sfreiberg/gotwilio"
	"google.golang.org/grpc/codes"
)

var util = objectify.New()

type Twilio struct {
	client *gotwilio.Twilio
}

func New(client *gotwilio.Twilio) *Twilio {
	return &Twilio{client}
}

func (t *Twilio) SendSMS(from, to *contacts.Contact, body, callback, appsid string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := t.client.SendSMS(from.Phone, to.Phone, body, callback, appsid)
	util.Debug("sent sms", "response", string(util.MarshalJSON(resp)))
	util.Debug("sent sms", "response", string(util.MarshalJSON(ex)))
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to send sms")
	}
	return resp, nil
}

func (t *Twilio) SendMMS(from, to *contacts.Contact, body, mediaurl, callback, appsid string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := t.client.SendMMS(from.Phone, to.Phone, body, mediaurl, callback, appsid)
	util.Debug("sent mms", "response", string(util.MarshalJSON(resp)))
	util.Debug("sent mms", "response", string(util.MarshalJSON(ex)))
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to send mms")
	}
	return resp, nil
}

func (t *Twilio) Call(from, to *contacts.Contact, callback string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := t.client.CallWithUrlCallbacks(from.Phone, to.Phone, gotwilio.NewCallbackParameters(callback))
	util.Debug("sent call", "response", string(util.MarshalJSON(resp)))
	util.Debug("sent call", "response", string(util.MarshalJSON(ex)))

	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to send call")
	}
	return resp, nil
}

func (t *Twilio) ProxyService(uniqueName, callback, outofsessionCallback, interceptCallback, geoMatchLevel, selectionBehavior string, defaultTTL int) (*gotwilio.ProxyService, error) {
	resp, ex, err := t.client.NewProxyService(gotwilio.ProxyServiceRequest{
		UniqueName:              uniqueName,
		CallbackURL:             callback,
		OutOfSessionCallbackURL: outofsessionCallback,
		InterceptCallbackURL:    interceptCallback,
		GeoMatchLevel:           geoMatchLevel,
		NumberSelectionBehavior: selectionBehavior,
		DefaultTtl:              defaultTTL,
	})
	util.Debug("sent proxy service", "response", string(util.MarshalJSON(resp)))
	util.Debug("sent proxy service", "response", string(util.MarshalJSON(ex)))
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to create proxy service")
	}

	return resp, nil
}

func (t *Twilio) VideoRoom() (*gotwilio.VideoResponse, error) {
	resp, ex, err := t.client.CreateVideoRoom(gotwilio.DefaultVideoRoomOptions)
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "failed to create proxy service")
	}

	return resp, nil
}
