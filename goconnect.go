package goconnect

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
	"firebase.google.com/go/auth"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

// GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client. It also carries an HTTP client and context.
type GoConnect struct {
	ctx   context.Context
	cli   *http.Client
	twil  *gotwilio.Twilio
	grid  *sendgrid.Client
	strip *client.API
	app   *firebase.App
}

// New Creates a new GoConnect from the provided http client, firebase credentials read from $PWN/credentials.json, and the following environmental
// variables: TWILIO_ACCOUNT TWILIO_TOKEN SENDGRID_TOKEN STRIPE_TOKEN
func New(cli *http.Client) *GoConnect {
	ctx := context.Background()
	if cli == nil {
		cli = http.DefaultClient
	}
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &GoConnect{
		twil:  gotwilio.NewTwilioClientCustomHTTP(os.Getenv("TWILIO_ACCOUNT"), os.Getenv("TWILIO_TOKEN"), cli),
		grid:  sendgrid.NewSendClient(os.Getenv("SENDGRID_TOKEN")),
		strip: client.New(os.Getenv("STRIPE_TOKEN"), stripe.NewBackends(cli)),
		cli:   cli,
		ctx : ctx,
		app: app,
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
