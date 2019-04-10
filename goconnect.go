package goconnect

import (
	"github.com/autom8ter/objectify"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	cli "github.com/stripe/stripe-go/client"
	"net/http"
	"os"
)

//	GoConnect holds the required configuration variables to create a GoConnect Instance. Use Init() to validate a GoConnect instance.
type GoConnect struct {
	TwilioAccount   string `validate:"required"`
	TwilioToken     string `validate:"required"`
	SendGridAccount string
	SendGridToken   string `validate:"required"`
	StripeAccount   string
	StripeToken     string `validate:"required"`
	SlackAccount    string
	SlackToken      string             `validate:"required"`
	util            *objectify.Handler `validate:"required"`
}

// New creates a new GoConnect Instance (no magic)
func New(twilioAccount string, twilioToken string, sendGridToken string, stripeAccount string, stripeToken string, slackToken string) *GoConnect {
	return &GoConnect{TwilioAccount: twilioAccount, TwilioToken: twilioToken, SendGridToken: sendGridToken, StripeAccount: stripeAccount, StripeToken: stripeToken, SlackToken: slackToken}
}

// NewFromFileEnv Initializes a gcp instance from service account credentials ref: https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
// and looks for Twilio, SendGrid, Stripe, and Slack credentials in environmental variables.
// vars: TWILIO_ACCOUNT, TWILIO_ACCOUNT, SENDGRID_ACCOUNT, SENDGRID_TOKEN, STRIPE_ACCOUNT, STRIPE_TOKEN, SLACK_ACCOUNT, SLACK_TOKEN
func NewFromFileEnv(file string) *GoConnect {
	return &GoConnect{
		TwilioAccount:   os.Getenv("TWILIO_ACCOUNT"),
		TwilioToken:     os.Getenv("TWILIO_ACCOUNT"),
		SendGridAccount: os.Getenv("SENDGRID_ACCOUNT"),
		SendGridToken:   os.Getenv("SENDGRID_TOKEN"),
		StripeAccount:   os.Getenv("STRIPE_ACCOUNT"),
		StripeToken:     os.Getenv("STRIPE_TOKEN"),
		SlackAccount:    os.Getenv("SLACK_ACCOUNT"),
		SlackToken:      os.Getenv("SLACK_TOKEN"),
		util:            objectify.Default(),
	}
}

// Init returns an error if any of the required fields are nil
func (g *GoConnect) Init() error {
	return g.util.Validate(g)
}

// Twilio returns an authenticated Twilio client
func (g *GoConnect) Twilio() *gotwilio.Twilio {
	g.util.PanicIfNil(g)
	return gotwilio.NewTwilioClient(g.TwilioAccount, g.TwilioToken)
}

// SendGrid returns an authenticated SendGrid client
func (g *GoConnect) SendGrid() *sendgrid.Client {
	g.util.PanicIfNil(g)
	return sendgrid.NewSendClient(g.SendGridToken)
}

//Stripe returns an authenticated Stripe client
func (g *GoConnect) Stripe(client *http.Client) *cli.API {
	g.util.PanicIfNil(g)
	return cli.New(g.SendGridToken, stripe.NewBackends(client))
}

//Util returns an objectify util tool ref:github.com/autom8ter/objectify
func (g *GoConnect) Util(client *http.Client) *objectify.Handler {
	return g.util
}
