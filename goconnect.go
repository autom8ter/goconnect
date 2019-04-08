package goconnect

import (
	"github.com/autom8ter/engine"
	"github.com/autom8ter/engine/config"
	"github.com/autom8ter/engine/driver"
	"github.com/autom8ter/gcloud"
	"github.com/autom8ter/objectify"
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
	StripeAccount   string
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
	tool.PanicIfNil(g)
	return tool.Validate(g)
}

// ToMap returns the GoConnect config as a map
func (g *GoConnect) ToMap() map[string]interface{} {
	tool.PanicIfNil(g)
	return tool.ToMap(g)
}

// Twilio returns an authenticated Twilio client
func (g *GoConnect) Twilio() *gotwilio.Twilio {
	tool.PanicIfNil(g)
	return gotwilio.NewTwilioClient(g.TwilioAccount, g.TwilioToken)
}

// SendGrid returns an authenticated SendGrid client
func (g *GoConnect) SendGrid() *sendgrid.Client {
	tool.PanicIfNil(g)
	return sendgrid.NewSendClient(g.SendGridToken)
}

//Stripe returns an authenticated Stripe client
func (g *GoConnect) Stripe(client *http.Client) *cli.API {
	tool.PanicIfNil(g)
	return cli.New(g.SendGridToken, stripe.NewBackends(client))
}

//Gcloud returns an authenticated GCP instance
func (g *GoConnect) Gcloud() *gcloud.GCP {
	tool.PanicIfNil(g.GCP)
	return g.GCP
}

// PluginFunc is a callback function that takes a GoConnect instance and returns a function that is used to create and register a grpc service.
// It is used in the GoConnect Serve() method.
type PluginFunc func(g *GoConnect) driver.PluginFunc

// Serve starts a grpc Engine ref:github.com/autom8ter/engine  server with a default middleware stack on the specified address with the provided pluugin functions.
func (g *GoConnect) Serve(addr string, fns ...PluginFunc) error {
	tool.PanicIfNil(g)
	plugs := []driver.Plugin{}
	for _, v := range fns {
		plugs = append(plugs, v(g))
	}
	return engine.New("tcp", addr, true).With(
		config.WithDefaultMiddlewares(),
		config.WithDefaultPlugins(),
		config.WithPlugins(plugs...),
	).Serve()
}
