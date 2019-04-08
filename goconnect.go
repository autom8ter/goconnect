package goconnect

import (
	"encoding/json"
	"fmt"
	"github.com/autom8ter/gcloud"
	"github.com/autom8ter/objectify"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	cli "github.com/stripe/stripe-go/client"
	"net/http"
)

var tool = objectify.New()

//	GoConnect holds the required configuration variables to create a GoConnect Instance. Use Init() to validate a GoConnect instance.
type GoConnect struct {
	GCP             *gcloud.GCP `validate:"required"`
	TwilioAccount   string      `validate:"required"`
	TwilioToken     string      `validate:"required"`
	SendGridAccount string
	SendGridToken   string `validate:"required"`
	StripeAccount   string
	StripeToken     string `validate:"required"`
	SlackAccount    string
	SlackToken      string `validate:"required"`
}

func New(g *gcloud.GCP, twilioAccount string, twilioToken string, sendGridToken string, stripeAccount string, stripeToken string, slackToken string) *GoConnect {
	return &GoConnect{GCP: g, TwilioAccount: twilioAccount, TwilioToken: twilioToken, SendGridToken: sendGridToken, StripeAccount: stripeAccount, StripeToken: stripeToken, SlackToken: slackToken}
}

func (g *GoConnect) Init() error {
	return tool.Validate(g)
}

func (g *GoConnect) Twilio() *gotwilio.Twilio {
	return gotwilio.NewTwilioClient(g.TwilioAccount, g.TwilioToken)
}

func (g *GoConnect) SendGrid() *sendgrid.Client {
	return sendgrid.NewSendClient(g.SendGridToken)
}

func (g *GoConnect) Stripe(client *http.Client) *cli.API {
	return cli.New(g.SendGridToken, stripe.NewBackends(client))
}
func (g *GoConnect) Gcloud() *gcloud.GCP {
	return g.GCP
}

// A HandlerFuncFunc is a GoConnect Callback function handler
type HandlerFunc func(g *GoConnect) error

// Execute runs the provided functions.
func (g *GoConnect) Execute(fns ...HandlerFunc) error {
	var err error
	for _, f := range fns {
		if newErr := f(g); newErr != nil {
			err = errors.Wrap(err, newErr.Error())
		}
	}
	return err
}

func toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
}
