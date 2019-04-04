package goconnect

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
	"github.com/autom8ter/engine/util"
	"github.com/autom8ter/gocrypt"
	"github.com/autom8ter/gosub"
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
	FirebaseCredsPath string `json:"firebase_creds_path"`
	TwilioAccount     string `json:"twilio_account"`
	TwilioToken       string `json:"twilio_token"`
	SendGridAccount   string `json:"sendgrid_account"`
	SendGridToken     string `json:"sendgrid_token"`
	StripeAccount     string `json:"stripe_account"`
	StripeToken       string `json:"stripe_token"`
	SlackAccount      string `json:"slack_account"`
	SlackToken        string `json:"slack_token"`
}

// GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client. It also carries an HTTP client and context.
type GoConnect struct {
	ctx   context.Context
	creds *Config
	cli   *http.Client
	twil  *gotwilio.Twilio
	grid  *sendgrid.Client
	strip *client.API
	app   *firebase.App
	chat  *slack.Client
	sub   *gosub.GoSub
	crypt *gocrypt.GoCrypt
}

// New Creates a new GoConnect from the provided http client and config
func New(cli *http.Client, c *Config) *GoConnect {
	ctx := context.Background()
	if cli == nil {
		cli = http.DefaultClient
	}
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(c.FirebaseCredsPath))
	if err != nil {
		log.Fatalln(err.Error())
	}
	su, err := gosub.New(ctx, "", option.WithCredentialsFile(c.FirebaseCredsPath))
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &GoConnect{
		creds: c,
		twil:  gotwilio.NewTwilioClientCustomHTTP(c.TwilioAccount, c.TwilioToken, cli),
		grid:  sendgrid.NewSendClient(c.SendGridToken),
		strip: client.New(c.StripeToken, stripe.NewBackends(cli)),
		cli:   cli,
		ctx:   ctx,
		app:   app,
		chat:  slack.New(c.SlackToken),
		sub:   su,
		crypt: gocrypt.NewGoCrypt(),
	}
}

// Stripe returns an authenticated Stripe client
func (g *GoConnect) Config() *Config {
	return g.creds
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

// Auth returns an authenticated Firebase Auth client
func (g *GoConnect) Auth() *auth.Client {
	a, err := g.app.Auth(g.ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return a
}

// Store returns an authenticated Firebase Storage client
func (g *GoConnect) Storage() *storage.Client {
	a, err := g.app.Storage(g.ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return a
}

// Firestore returns an authenticated Firebase Firestore client
func (g *GoConnect) Firestore() *firestore.Client {
	a, err := g.app.Firestore(g.ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return a
}

// Messaging returns an authenticated Firebase Messaging client
func (g *GoConnect) Messaging() *messaging.Client {
	a, err := g.app.Messaging(g.ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return a
}

// Database returns an authenticated Firebase Database client
func (g *GoConnect) Database(url string) *db.Client {
	a, err := g.app.DatabaseWithURL(g.ctx, url)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return a
}

// GoSub returns an authenticated GoSub client
func (g *GoConnect) GoSub() *gosub.GoSub {
	return g.sub
}

// GoSub returns a GoCrypt client
func (g *GoConnect) GoCrypt() *gocrypt.GoCrypt {
	return g.crypt
}

// Stringify formats an object and turns it into a JSON string
func (g *GoConnect) Stringify(obj interface{}) string {
	return util.ToPrettyJsonString(obj)
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
