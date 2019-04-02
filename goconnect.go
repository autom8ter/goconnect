package goconnect

import (
	"fmt"
	"github.com/autom8ter/goconnect/pkg/config"
	"github.com/autom8ter/goconnect/pkg/email"
	"github.com/autom8ter/goconnect/pkg/pay"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"net/http"
)

type Config struct {
	Client        *http.Client
	TwilioAccount string
	TwilioToken   string
	SendGridToken string
	StripeToken   string
}

type GoConnect struct {
	twil *gotwilio.Twilio
	grid *sendgrid.Client
}

func New(opts ...config.ConfigOption) *GoConnect {
	c := config.NewConfig(opts...)
	return &GoConnect{
		twil: gotwilio.NewTwilioClientCustomHTTP(c.Twilio.Account, c.Twilio.Token, c.Client),
		grid: sendgrid.NewSendClient(c.SendGrid.Token),
	}
}

func (g *GoConnect) NewCustomer(params *stripe.CustomerParams) (*stripe.Customer, error) {
	return customer.New(params)
}

func (g *GoConnect) ChargeCustomer(c *stripe.Customer, opts ...pay.ChargeOption) ([]*stripe.Charge, error) {
	params := pay.NewCharge(opts...)
	var charges = []*stripe.Charge{}
	i := charge.List(params)
	for i.Next() {
		charges = append(charges, i.Charge())
	}
	if err := i.Err(); err != nil {
		return charges, err
	}
	return charges, nil
}

func (g *GoConnect) SMS(to, from, body, callback, app string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := g.twil.SendSMS(from, to, body, callback, app)
	if err != nil {
		return resp, fmt.Errorf("exception: %s\nerror: %s\n", ex, err.Error())
	}
	return resp, nil
}

func (g *GoConnect) MMS(to, from, body, mediaURL string, callback, app string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := g.twil.SendMMS(from, to, body, mediaURL, callback, app)
	if err != nil {
		return resp, fmt.Errorf("exception: %s\nerror: %s\n", ex, err.Error())
	}
	return resp, nil
}

func (g *GoConnect) SMSCopilot(to, service, body, callback, app string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := g.twil.SendSMSWithCopilot(service, to, body, callback, app)
	if err != nil {
		return resp, fmt.Errorf("exception: %s\nerror: %s\n", ex, err.Error())
	}
	return resp, nil
}

func (g *GoConnect) Email(opts ...email.EmailOption) (*rest.Response, error) {
	return g.grid.Send(email.NewEmail(opts...))
}

func (g *GoConnect) Call(to, from, callback string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := g.twil.CallWithUrlCallbacks(from, to, gotwilio.NewCallbackParameters(callback))
	if err != nil {
		return resp, fmt.Errorf("exception: %s\nerror: %s\n", ex, err.Error())
	}
	return resp, nil
}

func (g *GoConnect) CallWithApp(to, from, appSid string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := g.twil.CallWithApplicationCallbacks(from, to, appSid)
	if err != nil {
		return resp, fmt.Errorf("exception: %s\nerror: %s\n", ex, err.Error())
	}
	return resp, nil
}
