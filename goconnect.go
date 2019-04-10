package goconnect

import (
	"context"
	"github.com/autom8ter/goconnect/hooks"
	"github.com/autom8ter/objectify"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/sub"
	"os"
)

type CustomerIndex int

const (
	ID CustomerIndex = iota
	EMAIL
	PHONE
)

type CallbackFunc func(customer2 *stripe.Customer) error

//	GoConnect holds the required configuration variables to create a GoConnect Instance. Use Init() to validate a GoConnect instance.
type GoConnect struct {
	twilio    *gotwilio.Twilio            `validate:"required"`
	grid      *sendgrid.Client            `validate:"required"`
	slck      *slack.Client               `validate:"required"`
	util      *objectify.Handler          `validate:"required"`
	hook      *hooks.SlackHook            `validate:"required"`
	customers map[string]*stripe.Customer `validate:"required"`
	cfg       *Config                     `validate:"required"`
}

type Config struct {
	TwilioAccount string `validate:"required"`
	TwilioKey     string `validate:"required"`
	SendgridKey   string `validate:"required"`
	SlackKey      string `validate:"required"`
	StripeKey     string `validate:"required"`
	Index         CustomerIndex
	EmailConfig   *EmailConfig `validate:"required"`
	LogConfig     *LogConfig   `validate:"required"`
}
type EmailConfig struct {
	Address string `validate:"required"`
	Name    string `validate:"required"`
}

type LogConfig struct {
	UserName string `validate:"required"`
	Channel  string `validate:"required"`
}

func New(cfg *Config) *GoConnect {
	util := objectify.Default()
	err := util.Validate(cfg)
	if err != nil {
		util.Fatalln(err.Error())
	}
	stripe.Key = cfg.StripeKey
	return &GoConnect{
		twilio:    gotwilio.NewTwilioClient(cfg.TwilioAccount, cfg.TwilioKey),
		grid:      sendgrid.NewSendClient(cfg.SendgridKey),
		slck:      slack.New(cfg.SlackKey),
		util:      util,
		customers: make(map[string]*stripe.Customer),
		cfg:       cfg,
	}
}

//TWILIO_ACCOUNT,  TWILIO_KEY,  SENDGRID_KEY,  SLACK_KEY,  STRIPE_KEY, EMAIL_ADDRESS, EMAIL_NAME, SLACK_LOG_USERNAME, SLACK_LOG_CHANNEL
func NewFromEnv(customerIndex CustomerIndex) *GoConnect {
	cfg := &Config{
		TwilioAccount: os.Getenv("TWILIO_ACCOUNT"),
		TwilioKey:     os.Getenv("TWILIO_KEY"),
		SendgridKey:   os.Getenv("SENDGRID_KEY"),
		SlackKey:      os.Getenv("SLACK_KEY"),
		StripeKey:     os.Getenv("STRIPE_KEY"),
		Index:         customerIndex,
		EmailConfig: &EmailConfig{
			Address: os.Getenv("EMAIL_ADDRESS"),
			Name:    os.Getenv("EMAIL_NAME"),
		},
		LogConfig: &LogConfig{
			UserName: os.Getenv("SLACK_LOG_USERNAME"),
			Channel:  os.Getenv("SLACK_LOG_CHANNEL"),
		},
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
		slck:      slack.New(cfg.SlackKey),
		util:      util,
		hook:      hooks.New(cfg.LogConfig.UserName, cfg.LogConfig.Channel, true),
		customers: make(map[string]*stripe.Customer),
		cfg:       cfg,
	}
}

// Init starts syncing the customer cache and validates the GoConnect instance
func (g *GoConnect) Init() error {
	go g.SyncCustomers()
	return g.util.Validate(g)
}

func (g *GoConnect) Config() *Config {
	return g.cfg
}

//Util returns an objectify util tool ref:github.com/autom8ter/objectify
func (g *GoConnect) Util() *objectify.Handler {
	return g.util
}

//Customers returns your current stripe customer cache
func (g *GoConnect) Customers() map[string]*stripe.Customer {
	return g.customers
}

func (g *GoConnect) SendEmail(name, address, subject, plain, html string) error {
	_, err := g.grid.Send(mail.NewSingleEmail(&mail.Email{
		Name:    g.cfg.EmailConfig.Name,
		Address: g.cfg.EmailConfig.Address,
	}, subject, &mail.Email{
		Name:    name,
		Address: address,
	}, plain, html))
	if err != nil {
		return err
	}
	return nil
}

//SendSMS sends an sms if mediaurl if empty, mms otherwise.
func (g *GoConnect) SendSMS(from, to, body, mediaUrl, callback, app string) (*gotwilio.SmsResponse, error) {
	if mediaUrl != "" {
		resp, ex, err := g.twilio.SendMMS(from, to, body, mediaUrl, callback, app)
		return resp, g.merge(ex, err)
	} else {
		resp, ex, err := g.twilio.SendSMS(from, to, body, callback, app)
		return resp, g.merge(ex, err)
	}
}

func (g *GoConnect) GetSMS(id string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := g.twilio.GetSMS(id)
	return resp, g.merge(ex, err)
}

func (g *GoConnect) GetCall(id string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := g.twilio.GetCall(id)
	return resp, g.merge(ex, err)
}

func (g *GoConnect) GetFax(id string) (*gotwilio.FaxResource, error) {
	resp, ex, err := g.twilio.GetFax(id)
	return resp, g.merge(ex, err)
}

func (g *GoConnect) GetVideoRoom(id string) (*gotwilio.VideoResponse, error) {
	resp, ex, err := g.twilio.GetVideoRoom(id)
	return resp, g.merge(ex, err)
}

//Call calls a number
func (g *GoConnect) SendCall(from, to, callback string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := g.twilio.CallWithUrlCallbacks(from, to, gotwilio.NewCallbackParameters(callback))
	return resp, g.merge(ex, err)
}

//Fax faxes a number
func (g *GoConnect) SendFax(to, from, mediaUrl, quality, callback string, storeMedia bool) (*gotwilio.FaxResource, error) {
	resp, ex, err := g.twilio.SendFax(to, from, mediaUrl, quality, callback, storeMedia)
	return resp, g.merge(ex, err)
}

//Fax faxes a number
func (g *GoConnect) CreateVideoRoom() (*gotwilio.VideoResponse, error) {
	resp, ex, err := g.twilio.CreateVideoRoom(gotwilio.DefaultVideoRoomOptions)
	return resp, g.merge(ex, err)
}

func (g *GoConnect) GetCustomer(key string) (*stripe.Customer, bool) {
	cust := g.customers[key]
	if cust != nil {
		return cust, true
	}
	return nil, false
}

func (g *GoConnect) SwitchIndex(typ CustomerIndex) {
	g.cfg.Index = typ
}

func (g *GoConnect) LogHook(ctx context.Context, author, icon, title string) error {
	return g.hook.PostLogEntry(ctx, g.slck, author, icon, title, g.util.Entry())
}

func (g *GoConnect) Hook(ctx context.Context, attachments ...slack.Attachment) error {
	return g.hook.PostAttachments(ctx, g.slck, attachments...)
}

func (g *GoConnect) SyncCustomers() {
	stripe.Key = g.cfg.StripeKey
	custList := customer.List(nil)
	c := custList.Customer()
	for {
		for custList.Next() {
			c = custList.Customer()
			switch g.cfg.Index {
			case EMAIL:
				g.customers[c.Email] = c
			case PHONE:
				g.customers[c.Shipping.Phone] = c
			default:
				g.customers[c.ID] = c
			}
		}
	}
}

func (g *GoConnect) GetSubscriptionFromCustomerEmail(email string) *stripe.Subscription {
	subList := sub.List(nil)
	s := subList.Subscription()
	for subList.Next() {
		s = subList.Subscription()
		if s.Customer.Email == email {
			return s
		}
	}
	return nil
}

func (g *GoConnect) GetSubscriptionFromCustomerPhone(phone string) *stripe.Subscription {
	subList := sub.List(nil)
	s := subList.Subscription()
	for subList.Next() {
		s = subList.Subscription()
		if s.Customer.Shipping.Phone == phone {
			return s
		}
	}
	return nil
}

func (g *GoConnect) GetSubscriptionFromCustomerID(id string) *stripe.Subscription {
	subList := sub.List(nil)
	s := subList.Subscription()
	for subList.Next() {
		s = subList.Subscription()
		if s.Customer.ID == id {
			return s
		}
	}
	return nil
}

func (g *GoConnect) SubscribeCustomer(key string, plan, cardnum, month, year, cvc string) (*stripe.Subscription, error) {
	if cust, ok := g.GetCustomer(key); ok {
		// create a subscription
		return sub.New(&stripe.SubscriptionParams{
			Customer: stripe.String(cust.ID),
			Plan:     stripe.String(plan),
			Card: &stripe.CardParams{
				Number:   stripe.String(cardnum),
				ExpMonth: stripe.String(month),
				ExpYear:  stripe.String(year),
				CVC:      stripe.String(cvc),
			},
		})
	} else {
		return nil, errors.New("customer not found: " + key)
	}
}

func (g *GoConnect) CancelSubscription(key string) error {
	if cust, ok := g.GetCustomer(key); ok {
		s := cust.Subscriptions.Data[0]
		_, err := sub.Cancel(s.ID, nil)
		if err != nil {
			return err
		}
	} else {
		return errors.New("customer not found: " + key)
	}
	return nil
}

func (g *GoConnect) CreateMonthlyPlan(amount int64, id, prodId, prodName, nickname string) (*stripe.Plan, error) {
	return plan.New(&stripe.PlanParams{
		Active:   stripe.Bool(true),
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("usd"),
		ID:       stripe.String(id),
		Interval: stripe.String("month"),
		Product: &stripe.PlanProductParams{
			Active: stripe.Bool(true),
			ID:     stripe.String(prodId),
			Name:   stripe.String(prodName),
		},
		Nickname: stripe.String(nickname),
	})
}

func (g *GoConnect) CreateYearlyPlan(amount int64, id, prodId, prodName, nickname string) (*stripe.Plan, error) {
	return plan.New(&stripe.PlanParams{
		Active:   stripe.Bool(true),
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("usd"),
		ID:       stripe.String(id),
		Interval: stripe.String("year"),
		Product: &stripe.PlanProductParams{
			Active: stripe.Bool(true),
			ID:     stripe.String(prodId),
			Name:   stripe.String(prodName),
		},
		Nickname: stripe.String(nickname),
	})
}

func (g *GoConnect) CreateCustomer(email, description, plan, name, phone string) (*stripe.Customer, error) {
	c, err := customer.New(&stripe.CustomerParams{
		Description: stripe.String(description),
		Email:       stripe.String(email),
		Plan:        stripe.String(plan),
		Shipping: &stripe.CustomerShippingDetailsParams{
			Address: &stripe.AddressParams{
				City:       stripe.String("N/A"),
				Country:    stripe.String("N/A"),
				Line1:      stripe.String("N/A"),
				Line2:      stripe.String("N/A"),
				PostalCode: stripe.String("N/A"),
				State:      stripe.String("N/A"),
			},
			Name:  stripe.String(name),
			Phone: stripe.String(phone),
		},
	})
	if err != nil {
		return nil, err
	}
	switch g.cfg.Index {
	case PHONE:
		g.customers[c.Shipping.Phone] = c
	case EMAIL:
		g.customers[c.Email] = c
	default:
		g.customers[c.ID] = c
	}
	return c, nil
}

func (g *GoConnect) CustomerKeys() []string {
	return g.util.Sort(g.customerKeys(g.customers))
}

func (g *GoConnect) CustomerExists(key string) bool {
	return g.util.Contains(g.CustomerKeys(), key)
}

func (g *GoConnect) CallBack(key string, funcs ...CallbackFunc) error {
	cust, ok := g.GetCustomer(key)
	if !ok {
		return errors.New("failed to find customer with key: " + key)
	}
	for i, f := range funcs {
		if err := f(cust); err != nil {
			return errors.Wrapf(err, "callback index: %v", i)
		}
	}
	return nil
}
