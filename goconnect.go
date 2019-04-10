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

var NOEXIST = func(key string) error {
	return errors.New("customer not found- key: " + key)
}

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
	Debug         bool
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
func NewFromEnv(customerIndex CustomerIndex, debug bool) *GoConnect {
	cfg := &Config{
		Debug:         debug,
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
	s := slack.New(cfg.SlackKey)
	if cfg.Debug {
		s.Debug()
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
		slck:      s,
		util:      util,
		hook:      hooks.New(cfg.LogConfig.UserName, cfg.LogConfig.Channel),
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

func (g *GoConnect) EmailCustomer(customerKey, subject, plain, html string) error {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return NOEXIST(customerKey)
	}
	_, err := g.grid.Send(mail.NewSingleEmail(&mail.Email{
		Name:    g.cfg.EmailConfig.Name,
		Address: g.cfg.EmailConfig.Address,
	}, subject, &mail.Email{
		Name:    cust.Shipping.Name,
		Address: cust.Email,
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

func (g *GoConnect) SMSCustomer(customerKey, from, body, mediaUrl, callback, app string) (*gotwilio.SmsResponse, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, NOEXIST(customerKey)
	}
	if mediaUrl != "" {
		resp, ex, err := g.twilio.SendMMS(from, cust.Shipping.Phone, body, mediaUrl, callback, app)
		return resp, g.merge(ex, err)
	} else {
		resp, ex, err := g.twilio.SendSMS(from, cust.Shipping.Phone, body, callback, app)
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

func (g *GoConnect) NewTwilioProxyService(name, callback, ofSessionCallback, interceptCallback, geoMatch, numSelectionBehavior string, defTTL int) (*gotwilio.ProxyService, error) {
	resp, ex, err := g.twilio.NewProxyService(gotwilio.ProxyServiceRequest{
		UniqueName:              name,
		CallbackURL:             callback,
		OutOfSessionCallbackURL: ofSessionCallback,
		InterceptCallbackURL:    interceptCallback,
		GeoMatchLevel:           geoMatch,
		NumberSelectionBehavior: numSelectionBehavior,
		DefaultTtl:              defTTL,
	})
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

//Call calls a number
func (g *GoConnect) CallCustomer(customerKey, from, callback string) (*gotwilio.VoiceResponse, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, NOEXIST(customerKey)
	}
	resp, ex, err := g.twilio.CallWithUrlCallbacks(from, cust.Shipping.Phone, gotwilio.NewCallbackParameters(callback))
	return resp, g.merge(ex, err)
}

//Fax faxes a number
func (g *GoConnect) SendFax(to, from, mediaUrl, quality, callback string, storeMedia bool) (*gotwilio.FaxResource, error) {
	resp, ex, err := g.twilio.SendFax(to, from, mediaUrl, quality, callback, storeMedia)
	return resp, g.merge(ex, err)
}

//Fax faxes a number
func (g *GoConnect) FaxCustomer(customerKey, from, mediaUrl, quality, callback string, storeMedia bool) (*gotwilio.FaxResource, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, NOEXIST(customerKey)
	}
	resp, ex, err := g.twilio.SendFax(cust.Shipping.Phone, from, mediaUrl, quality, callback, storeMedia)
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

func (g *GoConnect) CustomerMetadata(customerKey string) (map[string]string, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, NOEXIST(customerKey)
	}
	return cust.Metadata, nil
}

func (g *GoConnect) CustomerCard(customerKey string) (*stripe.Card, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, NOEXIST(customerKey)
	}
	return cust.DefaultSource.Card, nil
}

func (g *GoConnect) CustomerIsSubscribedToPlan(customerKey string, planFriendlyName string) bool {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return false
	}
	for _, s := range cust.Subscriptions.Data {
		if s.Plan.Nickname == planFriendlyName {
			return true
		}
	}
	return false
}

func (g *GoConnect) CustomerSubscriptions(customerKey string) ([]*stripe.Subscription, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, NOEXIST(customerKey)
	}
	return cust.Subscriptions.Data, nil
}

func (g *GoConnect) SubscribeCustomer(key string, plan, cardnum, month, year, cvc string) (*stripe.Subscription, error) {
	cust, ok := g.GetCustomer(key)
	if !ok {
		return nil, NOEXIST(key)
	}
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
}

func (g *GoConnect) CancelSubscription(key string, planName string) error {
	cust, ok := g.GetCustomer(key)
	if !ok {
		return NOEXIST(key)
	}
	for _, s := range cust.Subscriptions.Data {
		if s.Plan.Nickname == planName {
			_, err := sub.Cancel(s.ID, nil)
			if err != nil {
				return err
			}
		}
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
	if g.cfg.Index == EMAIL {
		_, ok := g.GetCustomer(email)
		if !ok {
			return nil, NOEXIST(email)
		}
	}
	if g.cfg.Index == PHONE {
		_, ok := g.GetCustomer(phone)
		if !ok {
			return nil, NOEXIST(phone)
		}
	}
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

func (g *GoConnect) CustomerCallBack(key string, funcs ...CallbackFunc) error {
	cust, ok := g.GetCustomer(key)
	if !ok {
		return NOEXIST(key)
	}
	for i, f := range funcs {
		if err := f(cust); err != nil {
			return errors.Wrapf(err, "callback index: %v", i)
		}
	}
	return nil
}

func (g *GoConnect) HandleSlackEvents(email string, funcs ...hooks.EventHandler) {
	hooks.EventLoop(g.slck, funcs...)
}

func (g *GoConnect) GetSlackThreadReplies(ctx context.Context, channel string, thread string) ([]slack.Message, error) {
	return g.slck.GetChannelRepliesContext(ctx, channel, thread)
}

func (g *GoConnect) GetSlackChannelHistory(ctx context.Context, channel, latest, oldest string, count int, inclusive bool) (*slack.History, error) {
	return g.slck.GetChannelHistoryContext(ctx, channel, slack.HistoryParameters{
		Latest:    latest,
		Oldest:    oldest,
		Count:     count,
		Inclusive: inclusive,
		Unreads:   true,
	})
}

func (g *GoConnect) GetUserByEmail(ctx context.Context, email string) (*slack.User, error) {
	return g.slck.GetUserByEmailContext(ctx, email)
}

func (g *GoConnect) UserIsAdmin(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsAdmin")
	}
	if usr.IsAdmin {
		return true, nil
	}

	return false, nil
}

func (g *GoConnect) UserIsPrimaryOwner(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsPrimaryOwner")
	}
	if usr.IsPrimaryOwner {
		return true, nil
	}

	return false, nil
}

func (g *GoConnect) UserIsOwner(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsOwner")
	}
	if usr.IsOwner {
		return true, nil
	}

	return false, nil
}

func (g *GoConnect) UserIsUltraRestricted(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsUltraRestricted")
	}
	if usr.IsUltraRestricted {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserIsAppUser(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsAppUser")
	}
	if usr.IsAppUser {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserIsBot(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsBot")
	}
	if usr.IsBot {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserIsStranger(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsStranger")
	}
	if usr.IsStranger {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserIsRestricted(ctx context.Context, email string) (bool, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsRestricted")
	}
	if usr.IsRestricted {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserPhoneNumber(ctx context.Context, email string) (string, error) {
	usr, err := g.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.Wrap(err, "goconnect.UserPhoneNumber")
	}
	return usr.Profile.Phone, nil
}

func (g *GoConnect) CallUser(ctx context.Context, email string, from string, callback string) (*gotwilio.VoiceResponse, error) {
	num, err := g.UserPhoneNumber(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "goconnect.CallUser")
	}
	return g.SendCall(from, num, callback)
}

func (g *GoConnect) SMSUser(ctx context.Context, email string, from string, body, mediaUrl string, callback, app string) (*gotwilio.SmsResponse, error) {
	num, err := g.UserPhoneNumber(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "goconnect.CallUser")
	}
	return g.SendSMS(from, num, body, mediaUrl, callback, app)
}

func (g *GoConnect) EmailUser(ctx context.Context, email, subject, string, plain, html string) error {
	usr, err := g.slck.GetUserByEmail(email)
	if err != nil {
		return errors.Wrap(err, "goconnect.EmailUser- Failed to get user by email")
	}
	return g.SendEmail(usr.Name, email, subject, plain, html)
}

func (g *GoConnect) AddChannelReminder(channelId string, text string, time string) (string, error) {
	rem, err := g.slck.AddChannelReminder(channelId, text, time)
	if err != nil {
		return "", errors.Wrap(err, "goconnect.AddChannelReminder")
	}
	return rem.ID, nil
}

func (g *GoConnect) AddUserReminder(userId string, text string, time string) (string, error) {
	rem, err := g.slck.AddUserReminder(userId, text, time)
	if err != nil {
		return "", errors.Wrap(err, "goconnect.AddUserReminder")
	}
	return rem.ID, nil
}

func (g *GoConnect) AddPin(ctx context.Context, text, channel, file, comment string) error {
	err := g.slck.AddPinContext(ctx, text, slack.ItemRef{
		Channel: channel,
		File:    channel,
		Comment: comment,
	})
	if err != nil {
		return errors.Wrap(err, "goconnect.AddPin")
	}
	return nil
}

func (g *GoConnect) AddStar(ctx context.Context, text, channel, file, comment string) error {
	err := g.slck.AddStarContext(ctx, text, slack.ItemRef{
		Channel: channel,
		File:    channel,
		Comment: comment,
	})
	if err != nil {
		return errors.Wrap(err, "goconnect.AddStar")
	}
	return nil
}

func (g *GoConnect) AddReaction(ctx context.Context, text, channel, file, comment string) error {
	err := g.slck.AddReactionContext(ctx, text, slack.ItemRef{
		Channel: channel,
		File:    channel,
		Comment: comment,
	})
	if err != nil {
		return errors.Wrap(err, "goconnect.AddReaction")
	}
	return nil
}
