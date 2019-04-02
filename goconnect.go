package goconnect

import (
	"context"
	"firebase.google.com/go"
	"github.com/autom8ter/goconnect/pkg/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

type GoConnect struct {
	ctx   context.Context
	debug bool
	cli   *http.Client
	twil  *gotwilio.Twilio
	grid  *sendgrid.Client
	strip *client.API
	app   *firebase.App
}

func New(opts ...config.Option) *GoConnect {
	c := config.NewConfig(opts...)
	return &GoConnect{
		twil:  gotwilio.NewTwilioClientCustomHTTP(c.Twilio.Account, c.Twilio.Token, c.Client),
		grid:  sendgrid.NewSendClient(c.SendGrid.Token),
		strip: client.New(c.Stripe.Token, stripe.NewBackends(c.Client)),
		cli:   c.Client,
		debug: c.Debug,
	}
}

func Default(cli *http.Client) *GoConnect {
	if cli == nil {
		cli = http.DefaultClient
	}
	c := config.NewConfig(config.FromEnv(cli))
	app, err := firebase.NewApp(c.Context, &firebase.Config{}, option.WithCredentialsFile(c.FirebaseCredsFile))
	if err != nil {
		log.Fatalln("failed to create firebase client", err.Error())
	}
	return &GoConnect{
		twil:  gotwilio.NewTwilioClientCustomHTTP(c.Twilio.Account, c.Twilio.Token, c.Client),
		grid:  sendgrid.NewSendClient(c.SendGrid.Token),
		strip: client.New(c.Stripe.Token, stripe.NewBackends(c.Client)),
		cli:   c.Client,
		debug: c.Debug,
		app:   app,
	}
}

func (g *GoConnect) Stripe() *client.API {
	return g.strip
}

func (g *GoConnect) Twilio() *gotwilio.Twilio {
	return g.twil
}

func (g *GoConnect) SendGrid() *sendgrid.Client {
	return g.grid
}

func (g *GoConnect) HTTP() *http.Client {
	return g.cli
}

func (g *GoConnect) Context() context.Context {
	return g.ctx
}

func (g *GoConnect) Firebase() *firebase.App {
	return g.app
}
