package goconnect

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
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
	SendGridToken     string `json:"sendgrid_token"`
	StripeToken       string `json:"stripe_token"`
}

// GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client. It also carries an HTTP client and context.
type GoConnect struct {
	ctx   context.Context
	cli   *http.Client
	twil  *gotwilio.Twilio
	grid  *sendgrid.Client
	strip *client.API
	app   *firebase.App
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
	return &GoConnect{
		twil:  gotwilio.NewTwilioClientCustomHTTP(c.TwilioAccount, c.TwilioToken, cli),
		grid:  sendgrid.NewSendClient(c.SendGridToken),
		strip: client.New(c.StripeToken, stripe.NewBackends(cli)),
		cli:   cli,
		ctx:   ctx,
		app:   app,
	}
}

// Stripe returns an authenticated Stripe client
func (g *GoConnect) Stripe() *client.API {
	return g.strip
}

// Twilio returns an authenticated Twilio client
func (g *GoConnect) Twilio() *gotwilio.Twilio {
	return g.twil
}

// Twilio returns an authenticated SendGrid client
func (g *GoConnect) SendGrid() *sendgrid.Client {
	return g.grid
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
