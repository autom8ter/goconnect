package goconnect

import (
	"context"
	"github.com/autom8ter/api/go/api"
	"github.com/autom8ter/goconnect/hooks"
	"github.com/autom8ter/slashsub"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/sub"
)

var CUSTOMERNOEXIST = func(key string) error {
	return errors.New("customer not found- key: " + key)
}

var USERNOEXIST = func(key string) error {
	return errors.New("customer not found- key: " + key)
}

//	GoConnect holds the required configuration variables to create a GoConnect Instance. Use Init() to validate a GoConnect instance.
type GoConnect struct {
	EmailAddress api.EmailAddress            `validate:"required"`
	PhoneNumber  string                      `validate:"required"`
	LogUsername  string                      `validate:"required"`
	LogChannel   string                      `validate:"required"`
	twilio       *gotwilio.Twilio            `validate:"required"`
	grid         *sendgrid.Client            `validate:"required"`
	slck         *slack.Client               `validate:"required"`
	hook         *hooks.SlackHook            `validate:"required"`
	customers    map[string]*stripe.Customer `validate:"required"`
	users        map[string]*slack.User      `validate:"required"`
	cfg          *api.Access                 `validate:"required"`
	slash        *slashsub.SlashSub          `validate:"required"`
}

func (g *GoConnect) Init(acc *api.Access) error {
	if acc == nil {
		acc = api.AccessFromEnv()
	}
	err := api.Util.Validate(acc)
	if err != nil {
		api.Util.Fatalln(err.Error())
	}

	stripe.Key = acc.StripeKey
	g.twilio = gotwilio.NewTwilioClient(acc.TwilioAccount, acc.TwilioKey)
	g.slck = slack.New(acc.SlackKey)
	g.customers = make(map[string]*stripe.Customer)
	g.cfg = acc
	g.grid = sendgrid.NewSendClient(acc.SendgridKey)
	g.slash, err = slashsub.New()
	if err != nil {
		return err
	}
	go g.SyncCustomers()
	return api.Util.Validate(g)
}

func (g *GoConnect) Access() *api.Access {
	return g.cfg
}

//Customers returns your current stripe customer cache
func (g *GoConnect) Customers() map[string]*stripe.Customer {
	return g.customers
}

func (g *GoConnect) Slash() *slashsub.SlashSub {
	return g.slash
}

func (g *GoConnect) SendEmail(r *api.RecipientEmail) error {
	resp, err := g.grid.Send(mail.NewSingleEmail(&mail.Email{
		Name:    g.EmailAddress.Name,
		Address: g.EmailAddress.Address,
	}, r.Subject, &mail.Email{
		Name:    r.To.Name,
		Address: r.To.Address,
	}, r.PlainText, r.Html))
	if err != nil {
		return errors.Wrap(err, string(api.Util.MarshalJSON(resp)))
	}
	return nil
}

func (g *GoConnect) EmailCustomer(customerKey, subject, plain, html string) error {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return CUSTOMERNOEXIST(customerKey)
	}

	return g.SendEmail(&api.RecipientEmail{
		To: &api.EmailAddress{
			Name:    cust.Shipping.Name,
			Address: cust.Email,
		},
		Subject:   subject,
		PlainText: plain,
		Html:      html,
	})
}

//SendSMS sends an sms if mediaurl if empty, mms otherwise.
func (g *GoConnect) SendSMS(s *api.SMS) (*gotwilio.SmsResponse, error) {
	if s.MediaUrl != "" {
		resp, ex, err := g.twilio.SendMMS(s.From, s.To, s.Body, s.MediaUrl, s.Callback, s.App)
		return resp, g.merge(ex, err)
	} else {
		resp, ex, err := g.twilio.SendSMS(s.From, s.To, s.Body, s.Callback, s.App)
		return resp, g.merge(ex, err)
	}
}

func (g *GoConnect) SMSCustomer(customerKey, from, body, mediaUrl, callback, app string) (*gotwilio.SmsResponse, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, CUSTOMERNOEXIST(customerKey)
	}
	return g.SendSMS(&api.SMS{
		To:       cust.Shipping.Phone,
		From:     from,
		Body:     body,
		MediaUrl: mediaUrl,
		Callback: callback,
		App:      app,
	})
}

func (g *GoConnect) GetSMS(id string) (*gotwilio.SmsResponse, error) {
	resp, ex, err := g.twilio.GetSMS(id)
	return resp, g.merge(ex, err)
}

func (g *GoConnect) GetCall(id string) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := g.twilio.GetCall(id)
	return resp, g.merge(ex, err)
}

func (g *GoConnect) NewTwilioProxyService(name, callback, outoffSessionCallback, interceptCallback, geoMatch, numSelectionBehavior string, defTTL int) (*gotwilio.ProxyService, error) {
	resp, ex, err := g.twilio.NewProxyService(gotwilio.ProxyServiceRequest{
		UniqueName:              name,
		CallbackURL:             callback,
		OutOfSessionCallbackURL: outoffSessionCallback,
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
func (g *GoConnect) SendCall(c *api.Call) (*gotwilio.VoiceResponse, error) {
	resp, ex, err := g.twilio.CallWithUrlCallbacks(c.From, c.To, gotwilio.NewCallbackParameters(c.Callback))
	return resp, g.merge(ex, err)
}

//Call calls a number
func (g *GoConnect) CallCustomer(customerKey, from, callback string) (*gotwilio.VoiceResponse, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, CUSTOMERNOEXIST(customerKey)
	}
	return g.SendCall(&api.Call{
		To:       cust.Shipping.Phone,
		From:     from,
		Callback: callback,
	})
}

//Fax faxes a number
func (g *GoConnect) SendFax(f *api.Fax) (*gotwilio.FaxResource, error) {
	resp, ex, err := g.twilio.SendFax(f.To, f.From, f.MediaUrl, f.Quality, f.Callback, f.StoreMedia)
	return resp, g.merge(ex, err)
}

//Fax faxes a number
func (g *GoConnect) FaxCustomer(customerKey, from, mediaUrl, quality, callback string, storeMedia bool) (*gotwilio.FaxResource, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, CUSTOMERNOEXIST(customerKey)
	}

	return g.SendFax(&api.Fax{
		To:         cust.Shipping.Phone,
		From:       from,
		Callback:   callback,
		MediaUrl:   mediaUrl,
		Quality:    quality,
		StoreMedia: storeMedia,
	})
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

func (g *GoConnect) LogHook(ctx context.Context, hook *api.LogHook) error {
	return g.hook.PostLogEntry(ctx, g.slck, hook.Author, hook.Icon, hook.Title, api.Util.Entry())
}

func (g *GoConnect) Hook(ctx context.Context, attachments ...api.Attachment) error {
	attachs := []slack.Attachment{}
	for _, a := range attachments {
		fields := []slack.AttachmentField{}
		for _, f := range a.Fields {
			fields = append(fields, slack.AttachmentField{
				Title: f.Title,
				Value: f.Value,
				Short: f.Short,
			})
		}

		attachs = append(attachs, slack.Attachment{
			Color:      a.Color,
			Fallback:   a.Fallback,
			CallbackID: a.CallbackId,
			ID:         int(a.Id),
			AuthorID:   a.AuthorId,
			AuthorName: a.AuthorName,
			AuthorLink: a.AuthorLink,
			AuthorIcon: a.AuthorIcon,
			Title:      a.Title,
			Pretext:    a.Pretext,
			Text:       a.Text,
			ImageURL:   a.ImageUrl,
			ThumbURL:   a.ThumbUrl,
			Fields:     fields,
		})
	}
	return g.hook.PostAttachments(ctx, g.slck, attachs...)
}

func (g *GoConnect) ActionHook(ctx context.Context, attachments ...api.Attachment) error {
	attachs := []slack.Attachment{}
	for _, a := range attachments {
		fields := []slack.AttachmentField{}
		for _, f := range a.Fields {
			fields = append(fields, slack.AttachmentField{
				Title: f.Title,
				Value: f.Value,
				Short: f.Short,
			})
		}

		attachs = append(attachs, slack.Attachment{
			Color:      a.Color,
			Fallback:   a.Fallback,
			CallbackID: a.CallbackId,
			ID:         int(a.Id),
			AuthorID:   a.AuthorId,
			AuthorName: a.AuthorName,
			AuthorLink: a.AuthorLink,
			AuthorIcon: a.AuthorIcon,
			Title:      a.Title,
			Pretext:    a.Pretext,
			Text:       a.Text,
			ImageURL:   a.ImageUrl,
			ThumbURL:   a.ThumbUrl,
			Fields:     fields,
		})
	}
	return g.hook.PostAttachments(ctx, g.slck, attachs...)
}
func (g *GoConnect) SyncCustomers() {
	stripe.Key = g.cfg.StripeKey
	custList := customer.List(nil)
	c := custList.Customer()
	for {
		for custList.Next() {
			g.customers[c.Email] = c
		}
	}
}

func (g *GoConnect) SyncUsers() {
	users, err := g.slck.GetUsers()
	if err != nil {
		api.Util.Fatalln(err.Error())
	}
	for _, usr := range users {
		g.users[usr.Profile.DisplayNameNormalized] = &usr
	}
}

func (g *GoConnect) CustomerMetadata(customerKey string) (map[string]string, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, CUSTOMERNOEXIST(customerKey)
	}
	return cust.Metadata, nil
}

func (g *GoConnect) CustomerCard(customerKey string) (*stripe.Card, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, CUSTOMERNOEXIST(customerKey)
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
		return nil, CUSTOMERNOEXIST(customerKey)
	}
	return cust.Subscriptions.Data, nil
}

func (g *GoConnect) SubscribeCustomer(key string, plan, cardnum, month, year, cvc string) (*stripe.Subscription, error) {
	cust, ok := g.GetCustomer(key)
	if !ok {
		return nil, CUSTOMERNOEXIST(key)
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
		return CUSTOMERNOEXIST(key)
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
	_, ok := g.GetCustomer(email)
	if !ok {
		return nil, CUSTOMERNOEXIST(email)
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
	g.customers[c.Email] = c
	return c, nil
}

func (g *GoConnect) CustomerKeys() []string {
	return api.Util.Sort(g.customerKeys(g.customers))
}

func (g *GoConnect) CustomerExists(key string) bool {
	return api.Util.Contains(g.CustomerKeys(), key)
}

func (g *GoConnect) HandleSlackEvents(funcs ...hooks.EventHandler) {
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

func (g *GoConnect) UserIsAdmin(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsAdmin")
	}
	if usr.IsAdmin {
		return true, nil
	}

	return false, nil
}

func (g *GoConnect) UserIsPrimaryOwner(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsPrimaryOwner")
	}
	if usr.IsPrimaryOwner {
		return true, nil
	}

	return false, nil
}

func (g *GoConnect) UserIsOwner(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsOwner")
	}
	if usr.IsOwner {
		return true, nil
	}

	return false, nil
}

func (g *GoConnect) UserIsUltraRestricted(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsUltraRestricted")
	}
	if usr.IsUltraRestricted {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserIsAppUser(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsAppUser")
	}
	if usr.IsAppUser {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserIsBot(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsBot")
	}
	if usr.IsBot {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserKeys() []string {
	return api.Util.Sort(g.userKeys(g.users))
}

func (g *GoConnect) UserExists(key string) bool {
	return api.Util.Contains(g.UserKeys(), key)
}

func (g *GoConnect) UserIsStranger(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsStranger")
	}
	if usr.IsStranger {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserIsRestricted(ctx context.Context, name string) (bool, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return false, errors.Wrap(err, "goconnect.UserIsRestricted")
	}
	if usr.IsRestricted {
		return true, nil
	}
	return false, nil
}

func (g *GoConnect) UserPhoneNumber(ctx context.Context, name string) (string, error) {
	usr, err := g.GetUser(name)
	if err != nil {
		return "", errors.Wrap(err, "goconnect.UserPhoneNumber")
	}
	return usr.Profile.Phone, nil
}

func (g *GoConnect) CallUser(ctx context.Context, r *api.CallRequest) (*gotwilio.VoiceResponse, error) {
	usr, err := g.GetUser(r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "goconnect.CallUser")
	}
	return g.SendCall(&api.Call{
		To:       usr.Profile.Phone,
		From:     g.PhoneNumber,
		Callback: r.CallbackUrl,
	})
}

func (g *GoConnect) SMSUser(ctx context.Context, r *api.SMSRequest) (*gotwilio.SmsResponse, error) {
	num, err := g.UserPhoneNumber(ctx, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "goconnect.SMSUser")
	}
	return g.SendSMS(&api.SMS{
		To:   num,
		From: g.PhoneNumber,
		Body: r.Body,
	})
}

func (g *GoConnect) MMSUser(ctx context.Context, r *api.MMSRequest) (*gotwilio.SmsResponse, error) {
	num, err := g.UserPhoneNumber(ctx, r.Sms.Id)
	if err != nil {
		return nil, errors.Wrap(err, "goconnect.MMSUser")
	}
	return g.SendSMS(&api.SMS{
		To:   num,
		From: g.PhoneNumber,
		Body: r.Sms.Body,
	})
}

func (g *GoConnect) GetUser(name string) (*slack.User, error) {
	user := g.users[name]
	if user == nil {
		g.SyncUsers()
	}
	if user == nil {
		return nil, USERNOEXIST(name)
	}
	return user, nil
}

func (g *GoConnect) EmailUser(ctx context.Context, r *api.EmailRequest) error {
	usr, err := g.GetUser(r.Id)
	if err != nil {
		return err
	}
	err = g.SendEmail(&api.RecipientEmail{
		To: &api.EmailAddress{
			Name:    usr.Name,
			Address: usr.Profile.Email,
		},
		Subject:   r.Subject,
		PlainText: r.PlainText,
		Html:      r.HtmlAlt,
	})
	return err
}

func (g *GoConnect) AddChannelReminder(r *api.ChannelReminder) (string, error) {
	rem, err := g.slck.AddChannelReminder(r.ChannelId, r.Text, r.Time)
	if err != nil {
		return "", errors.Wrap(err, "goconnect.AddChannelReminder")
	}
	return rem.ID, nil
}

func (g *GoConnect) AddUserReminder(r *api.UserReminder) (*slack.Reminder, error) {
	rem, err := g.slck.AddUserReminder(r.Id, r.Text, r.Time)
	if err != nil {
		return nil, errors.Wrap(err, "goconnect.AddUserReminder")
	}
	return rem, nil
}

func (g *GoConnect) AddPin(ctx context.Context, p *api.Pin) error {
	err := g.slck.AddPinContext(ctx, p.Text, slack.ItemRef{
		Channel: p.Item.Channel,
		File:    p.Item.File,
		Comment: p.Item.Comment,
	})
	if err != nil {
		return errors.Wrap(err, "goconnect.AddPin")
	}
	return nil
}

func (g *GoConnect) AddStar(ctx context.Context, star *api.Star) error {
	err := g.slck.AddStarContext(ctx, star.Text, slack.ItemRef{
		Channel: star.Item.Channel,
		File:    star.Item.File,
		Comment: star.Item.Comment,
	})
	if err != nil {
		return errors.Wrap(err, "goconnect.AddStar")
	}
	return nil
}

func (g *GoConnect) AddReaction(ctx context.Context, r *api.UserReminder) error {
	err := g.slck.AddReactionContext(ctx, r.Text, slack.ItemRef{
		Channel: r.Item.Channel,
		File:    r.Item.File,
		Comment: r.Item.Comment,
	})
	if err != nil {
		return errors.Wrap(err, "goconnect.AddReaction")
	}
	return nil
}
