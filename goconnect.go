package goconnect

import (
	"github.com/autom8ter/gcloud"
	"github.com/autom8ter/objectify"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	cli "github.com/stripe/stripe-go/client"
	"google.golang.org/api/option"
	"net/http"
	"os"
)

var tool = objectify.New()

//	GoConnect holds the required configuration variables to create a GoConnect Instance. Use Init() to validate a GoConnect instance.
type GoConnect struct {
	GCP             *gcloud.GCP `validate:"required"`
	TwilioAccount   string      `validate:"required"`
	TwilioToken     string      `validate:"required"`
	SendGridAccount string
	SendGridToken   string `validate:"required"`
	StripeAccount   string `validate:"required"`
	StripeToken     string `validate:"required"`
	SlackAccount    string
	SlackToken      string `validate:"required"`
}

// New creates a new GoConnect Instance (no magic)
func New(twilioAccount string, twilioToken string, sendGridToken string, stripeAccount string, stripeToken string, slackToken string, opts ...option.ClientOption) *GoConnect {
	return &GoConnect{GCP: gcloud.NewGCP(opts...), TwilioAccount: twilioAccount, TwilioToken: twilioToken, SendGridToken: sendGridToken, StripeAccount: stripeAccount, StripeToken: stripeToken, SlackToken: slackToken}
}

// NewFromFileEnv Initializes a gcp instance from service account credentials ref: https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
// and looks for Twilio, SendGrid, Stripe, and Slack credentials in environmental variables.
// vars: TWILIO_ACCOUNT, TWILIO_ACCOUNT, SENDGRID_ACCOUNT, SENDGRID_TOKEN, STRIPE_ACCOUNT, STRIPE_TOKEN, SLACK_ACCOUNT, SLACK_TOKEN
func NewFromFileEnv(file string) *GoConnect {
	return &GoConnect{
		GCP:             gcloud.NewGCP(option.WithCredentialsFile(file)),
		TwilioAccount:   os.Getenv("TWILIO_ACCOUNT"),
		TwilioToken:     os.Getenv("TWILIO_ACCOUNT"),
		SendGridAccount: os.Getenv("SENDGRID_ACCOUNT"),
		SendGridToken:   os.Getenv("SENDGRID_TOKEN"),
		StripeAccount:   os.Getenv("STRIPE_ACCOUNT"),
		StripeToken:     os.Getenv("STRIPE_TOKEN"),
		SlackAccount:    os.Getenv("SLACK_ACCOUNT"),
		SlackToken:      os.Getenv("SLACK_TOKEN"),
	}
}

// Init returns an error if any of the required fields are nil
func (g *GoConnect) Init() error {
	return tool.Validate(g)
}

// ToMap returns the GoConnect config as a map
func (g *GoConnect) ToMap() map[string]interface{} {
	return tool.ToMap(g)
}

// Twilio returns an authenticated Twilio client
func (g *GoConnect) Twilio() *gotwilio.Twilio {
	return gotwilio.NewTwilioClient(g.TwilioAccount, g.TwilioToken)
}

// SendGrid returns an authenticated SendGrid client
func (g *GoConnect) SendGrid() *sendgrid.Client {
	return sendgrid.NewSendClient(g.SendGridToken)
}

//Stripe returns an authenticated Stripe client
func (g *GoConnect) Stripe(client *http.Client) *cli.API {
	return cli.New(g.SendGridToken, stripe.NewBackends(client))
}

//Gcloud returns an authenticated GCP instance
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
