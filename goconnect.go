package goconnect

import (
	"context"
	"github.com/autom8ter/gcloud"
	"github.com/autom8ter/gocrypt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

//Config holds the required configuration variables to create a GoConnect Instance
type Config struct {
	GCPCredsPath    string `json:"firebase_creds_path"`
	TwilioAccount   string `json:"twilio_account"`
	TwilioToken     string `json:"twilio_token"`
	SendGridAccount string `json:"sendgrid_account"`
	SendGridToken   string `json:"sendgrid_token"`
	StripeAccount   string `json:"stripe_account"`
	StripeToken     string `json:"stripe_token"`
	SlackAccount    string `json:"slack_account"`
	SlackToken      string `json:"slack_token"`
}

// GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client. It also carries an HTTP client and context.
type GoConnect struct {
	ctx   context.Context
	cfg   *Config
	twil  *gotwilio.Twilio
	grid  *sendgrid.Client
	strip *client.API
	chat  *slack.Client
	gcp   *gcloud.GCP
	crypt *gocrypt.GoCrypt
	cli   *http.Client
}

// New Creates a new GoConnect from the provided http client and config
func New(cli *http.Client, c *Config) *GoConnect {
	ctx := context.Background()
	if cli == nil {
		cli = http.DefaultClient
	}

	gcp, err := gcloud.New(ctx, option.WithCredentialsFile(c.GCPCredsPath))
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &GoConnect{
		cfg:   c,
		twil:  gotwilio.NewTwilioClientCustomHTTP(c.TwilioAccount, c.TwilioToken, cli),
		grid:  sendgrid.NewSendClient(c.SendGridToken),
		strip: client.New(c.StripeToken, stripe.NewBackends(cli)),
		ctx:   ctx,
		chat:  slack.New(c.SlackToken),
		crypt: gocrypt.NewGoCrypt(),
		gcp:   gcp,
	}
}

// GCP returns a gcloud.GCP instance
func (g *GoConnect) GCP() *gcloud.GCP {
	return g.gcp
}

// Config returns the config used to create the GoConnect instance
func (g *GoConnect) Config() *Config {
	return g.cfg
}

// Stripe returns an authenticated Stripe client
func (g *GoConnect) Stripe() *client.API {
	return g.strip
}

// Twilio returns an authenticated Twilio client
func (g *GoConnect) Twilio() *gotwilio.Twilio {
	return g.twil
}

// SendGrid returns an authenticated SendGrid client
func (g *GoConnect) SendGrid() *sendgrid.Client {
	return g.grid
}

// Slack returns an authenticated Slack client
func (g *GoConnect) Slack() *slack.Client {
	return g.chat
}

// Twilio returns an HTTP client
func (g *GoConnect) HTTP() *http.Client {
	return g.cli
}

// GoSub returns a GoCrypt client
func (g *GoConnect) GoCrypt() *gocrypt.GoCrypt {
	return g.crypt
}

type Func func(g *GoConnect) error

func (g *GoConnect) Execute(fns ...Func) error {
	var err error
	for _, f := range fns {
		if newErr := f(g); newErr != nil {
			err = errors.Wrap(err, newErr.Error())
		}
	}
	return err
}
