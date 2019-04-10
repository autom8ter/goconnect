package goconnect

import (
	"github.com/autom8ter/objectify"
	"github.com/nlopes/slack"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"os"
	"time"
)

type CustomerIndex int

const (
	ID CustomerIndex = iota
	EMAIL
	PHONE
)

//	GoConnect holds the required configuration variables to create a GoConnect Instance. Use Init() to validate a GoConnect instance.
type GoConnect struct {
	twilio    *gotwilio.Twilio            `validate:"required"`
	grid      *sendgrid.Client            `validate:"required"`
	slck      *slack.Client               `validate:"required"`
	util      *objectify.Handler          `validate:"required"`
	customers map[string]*stripe.Customer `validate:"required"`
	cfg       *Config                     `validate:"required"`
}

type Config struct {
	TwilioAccount string        `validate:"required"`
	TwilioKey     string        `validate:"required"`
	SendgridKey   string        `validate:"required"`
	SlackKey      string        `validate:"required"`
	StripeKey     string        `validate:"required"`
	Index         CustomerIndex `validate:"required"`
	SyncFrequency string        `validate:"required"`
}

func New(cfg *Config) *GoConnect {
	util := objectify.Default()
	err := util.Validate(cfg)
	if err != nil {
		util.Fatalln(err.Error())
	}
	stripe.Key = cfg.StripeKey
	customers := make(map[string]*stripe.Customer)
	i := customer.List(nil)
	c := i.Customer()
	for i.Next() {
		switch cfg.Index {
		case EMAIL:
			customers[c.Email] = c
		case PHONE:
			customers[c.Shipping.Phone] = c
		default:
			customers[c.ID] = c
		}
	}

	return &GoConnect{
		twilio:    gotwilio.NewTwilioClient(cfg.TwilioAccount, cfg.TwilioKey),
		grid:      sendgrid.NewSendClient(cfg.SendgridKey),
		slck:      slack.New(cfg.SlackKey),
		util:      util,
		customers: customers,
		cfg:       cfg,
	}
}

// NewFromFileEnv Initializes a gcp instance from service account credentials ref: https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
// and looks for Twilio, SendGrid, Stripe, and Slack credentials in environmental variables.
// vars: TWILIO_ACCOUNT, TWILIO_KEY, SENDGRID_KEY, SYNC_FREQUENCY, STRIPE_KEY, STRIPE_TOKEN, SLACK_KEY
func NewFromEnv() *GoConnect {
	cfg := &Config{
		TwilioAccount: os.Getenv("TWILIO_ACCOUNT"),
		TwilioKey:     os.Getenv("TWILIO_KEY"),
		SendgridKey:   os.Getenv("SENDGRID_KEY"),
		SlackKey:      os.Getenv("SLACK_KEY"),
		StripeKey:     os.Getenv("STRIPE_KEY"),
		Index:         EMAIL,
		SyncFrequency: os.Getenv("SYNC_FREQUENCY"),
	}
	util := objectify.Default()
	err := util.Validate(cfg)
	if err != nil {
		util.Fatalln(err.Error())
	}
	stripe.Key = cfg.StripeKey
	return &GoConnect{
		twilio:    gotwilio.NewTwilioClient(cfg.TwilioAccount, cfg.TwilioKey),
		grid:      sendgrid.NewSendClient(cfg.SendgridKey),
		slck:      nil,
		util:      util,
		customers: make(map[string]*stripe.Customer),
		cfg:       cfg,
	}
}

// Init starts syncing the customer cache and validates the GoConnect instance
func (g *GoConnect) Init() error {
	freq, err := time.ParseDuration(g.cfg.SyncFrequency)
	if err != nil {
		return err
	}
	if freq != 0 {
		g.sync(freq)
	} else {
		g.sync(1 * time.Minute)
	}
	return g.util.Validate(g)
}

//Util returns an objectify util tool ref:github.com/autom8ter/objectify
func (g *GoConnect) Util() *objectify.Handler {
	return g.util
}

//Customers returns your current stripe customer cache
func (g *GoConnect) Customers() map[string]*stripe.Customer {
	return g.customers
}

//SendSMS sends an sms if mediaurl if empty, mms otherwise.
func (g *GoConnect) SMS(from, to, body, mediaUrl, callback, app string) (*gotwilio.SmsResponse, error) {
	if mediaUrl != "" {
		resp, ex, err := g.twilio.SendMMS(from, to, body, mediaUrl, callback, app)
		return resp, g.merge(ex, err)
	} else {
		resp, ex, err := g.twilio.SendSMS(from, to, body, callback, app)
		return resp, g.merge(ex, err)
	}
}

//Call calls a number
func (g *GoConnect) Call(from, to, callback string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := g.twilio.CallWithUrlCallbacks(from, to, gotwilio.NewCallbackParameters(callback))
	return resp, g.merge(ex, err)

}

func (g *GoConnect) merge(ex *gotwilio.Exception, err error) error {
	if err != nil && ex != nil {
		return g.Util().WrapErr(err, string(g.Util().MarshalJSON(ex)))
	}
	if err != nil {
		return err
	}
	return nil
}

func (g *GoConnect) GetCustomer(key string) *stripe.Customer {
	return g.customers[key]
}

func (g *GoConnect) SwitchIndex(typ CustomerIndex) {
	g.cfg.Index = typ
}

func (g *GoConnect) sync(frequency time.Duration) {
	stripe.Key = g.cfg.StripeKey
	for {
		i := customer.List(nil)
		c := i.Customer()
		for i.Next() {
			switch g.cfg.Index {
			case EMAIL:
				g.customers[c.Email] = c
			case PHONE:
				g.customers[c.Shipping.Phone] = c
			default:
				g.customers[c.ID] = c
			}
		}
		time.Sleep(frequency)
	}
}
