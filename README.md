# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

```go
var NOEXIST = func(key string) error {
	return errors.New("customer not found- key: " + key)
}
```

#### type CallbackFunc

```go
type CallbackFunc func(customer2 *stripe.Customer) error
```


#### type Config

```go
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
```


#### type CustomerIndex

```go
type CustomerIndex int
```


```go
const (
	ID CustomerIndex = iota
	EMAIL
	PHONE
)
```

#### type EmailConfig

```go
type EmailConfig struct {
	Address string `validate:"required"`
	Name    string `validate:"required"`
}
```


#### type GoConnect

```go
type GoConnect struct {
}
```

GoConnect holds the required configuration variables to create a GoConnect
Instance. Use Init() to validate a GoConnect instance.

#### func  New

```go
func New(cfg *Config) *GoConnect
```

#### func  NewFromEnv

```go
func NewFromEnv(customerIndex CustomerIndex, debug bool) *GoConnect
```
TWILIO_ACCOUNT, TWILIO_KEY, SENDGRID_KEY, SLACK_KEY, STRIPE_KEY, EMAIL_ADDRESS,
EMAIL_NAME, SLACK_LOG_USERNAME, SLACK_LOG_CHANNEL

#### func (*GoConnect) AddChannelReminder

```go
func (g *GoConnect) AddChannelReminder(channelId string, text string, time string) (string, error)
```

#### func (*GoConnect) AddPin

```go
func (g *GoConnect) AddPin(ctx context.Context, text, channel, file, comment string) error
```

#### func (*GoConnect) AddReaction

```go
func (g *GoConnect) AddReaction(ctx context.Context, text, channel, file, comment string) error
```

#### func (*GoConnect) AddStar

```go
func (g *GoConnect) AddStar(ctx context.Context, text, channel, file, comment string) error
```

#### func (*GoConnect) AddUserReminder

```go
func (g *GoConnect) AddUserReminder(userId string, text string, time string) (string, error)
```

#### func (*GoConnect) CallCustomer

```go
func (g *GoConnect) CallCustomer(customerKey, from, callback string) (*gotwilio.VoiceResponse, error)
```
Call calls a number

#### func (*GoConnect) CallUser

```go
func (g *GoConnect) CallUser(ctx context.Context, email string, from string, callback string) (*gotwilio.VoiceResponse, error)
```

#### func (*GoConnect) CancelSubscription

```go
func (g *GoConnect) CancelSubscription(key string, planName string) error
```

#### func (*GoConnect) Config

```go
func (g *GoConnect) Config() *Config
```

#### func (*GoConnect) CreateCustomer

```go
func (g *GoConnect) CreateCustomer(email, description, plan, name, phone string) (*stripe.Customer, error)
```

#### func (*GoConnect) CreateMonthlyPlan

```go
func (g *GoConnect) CreateMonthlyPlan(amount int64, id, prodId, prodName, nickname string) (*stripe.Plan, error)
```

#### func (*GoConnect) CreateVideoRoom

```go
func (g *GoConnect) CreateVideoRoom() (*gotwilio.VideoResponse, error)
```
Fax faxes a number

#### func (*GoConnect) CreateYearlyPlan

```go
func (g *GoConnect) CreateYearlyPlan(amount int64, id, prodId, prodName, nickname string) (*stripe.Plan, error)
```

#### func (*GoConnect) CustomerCallBack

```go
func (g *GoConnect) CustomerCallBack(key string, funcs ...CallbackFunc) error
```

#### func (*GoConnect) CustomerCard

```go
func (g *GoConnect) CustomerCard(customerKey string) (*stripe.Card, error)
```

#### func (*GoConnect) CustomerExists

```go
func (g *GoConnect) CustomerExists(key string) bool
```

#### func (*GoConnect) CustomerIsSubscribedToPlan

```go
func (g *GoConnect) CustomerIsSubscribedToPlan(customerKey string, planFriendlyName string) bool
```

#### func (*GoConnect) CustomerKeys

```go
func (g *GoConnect) CustomerKeys() []string
```

#### func (*GoConnect) CustomerMetadata

```go
func (g *GoConnect) CustomerMetadata(customerKey string) (map[string]string, error)
```

#### func (*GoConnect) CustomerSubscriptions

```go
func (g *GoConnect) CustomerSubscriptions(customerKey string) ([]*stripe.Subscription, error)
```

#### func (*GoConnect) Customers

```go
func (g *GoConnect) Customers() map[string]*stripe.Customer
```
Customers returns your current stripe customer cache

#### func (*GoConnect) EmailCustomer

```go
func (g *GoConnect) EmailCustomer(customerKey, subject, plain, html string) error
```

#### func (*GoConnect) EmailUser

```go
func (g *GoConnect) EmailUser(ctx context.Context, email, subject, string, plain, html string) error
```

#### func (*GoConnect) FaxCustomer

```go
func (g *GoConnect) FaxCustomer(customerKey, from, mediaUrl, quality, callback string, storeMedia bool) (*gotwilio.FaxResource, error)
```
Fax faxes a number

#### func (*GoConnect) GetCall

```go
func (g *GoConnect) GetCall(id string) (*gotwilio.VoiceResponse, error)
```

#### func (*GoConnect) GetCustomer

```go
func (g *GoConnect) GetCustomer(key string) (*stripe.Customer, bool)
```

#### func (*GoConnect) GetFax

```go
func (g *GoConnect) GetFax(id string) (*gotwilio.FaxResource, error)
```

#### func (*GoConnect) GetSMS

```go
func (g *GoConnect) GetSMS(id string) (*gotwilio.SmsResponse, error)
```

#### func (*GoConnect) GetSlackChannelHistory

```go
func (g *GoConnect) GetSlackChannelHistory(ctx context.Context, channel, latest, oldest string, count int, inclusive bool) (*slack.History, error)
```

#### func (*GoConnect) GetSlackThreadReplies

```go
func (g *GoConnect) GetSlackThreadReplies(ctx context.Context, channel string, thread string) ([]slack.Message, error)
```

#### func (*GoConnect) GetUserByEmail

```go
func (g *GoConnect) GetUserByEmail(ctx context.Context, email string) (*slack.User, error)
```

#### func (*GoConnect) GetVideoRoom

```go
func (g *GoConnect) GetVideoRoom(id string) (*gotwilio.VideoResponse, error)
```

#### func (*GoConnect) HandleSlackEvents

```go
func (g *GoConnect) HandleSlackEvents(email string, funcs ...hooks.EventHandler)
```

#### func (*GoConnect) Hook

```go
func (g *GoConnect) Hook(ctx context.Context, attachments ...slack.Attachment) error
```

#### func (*GoConnect) Init

```go
func (g *GoConnect) Init() error
```
Init starts syncing the customer cache and validates the GoConnect instance

#### func (*GoConnect) LogHook

```go
func (g *GoConnect) LogHook(ctx context.Context, author, icon, title string) error
```

#### func (*GoConnect) NewTwilioProxyService

```go
func (g *GoConnect) NewTwilioProxyService(name, callback, ofSessionCallback, interceptCallback, geoMatch, numSelectionBehavior string, defTTL int) (*gotwilio.ProxyService, error)
```

#### func (*GoConnect) SMSCustomer

```go
func (g *GoConnect) SMSCustomer(customerKey, from, body, mediaUrl, callback, app string) (*gotwilio.SmsResponse, error)
```

#### func (*GoConnect) SMSUser

```go
func (g *GoConnect) SMSUser(ctx context.Context, email string, from string, body, mediaUrl string, callback, app string) (*gotwilio.SmsResponse, error)
```

#### func (*GoConnect) SendCall

```go
func (g *GoConnect) SendCall(from, to, callback string) (*gotwilio.VoiceResponse, error)
```
Call calls a number

#### func (*GoConnect) SendEmail

```go
func (g *GoConnect) SendEmail(name, address, subject, plain, html string) error
```

#### func (*GoConnect) SendFax

```go
func (g *GoConnect) SendFax(to, from, mediaUrl, quality, callback string, storeMedia bool) (*gotwilio.FaxResource, error)
```
Fax faxes a number

#### func (*GoConnect) SendSMS

```go
func (g *GoConnect) SendSMS(from, to, body, mediaUrl, callback, app string) (*gotwilio.SmsResponse, error)
```
SendSMS sends an sms if mediaurl if empty, mms otherwise.

#### func (*GoConnect) Serve

```go
func (g *GoConnect) Serve(addr string, plugs ...Plugin) error
```

#### func (*GoConnect) SubscribeCustomer

```go
func (g *GoConnect) SubscribeCustomer(key string, plan, cardnum, month, year, cvc string) (*stripe.Subscription, error)
```

#### func (*GoConnect) SwitchIndex

```go
func (g *GoConnect) SwitchIndex(typ CustomerIndex)
```

#### func (*GoConnect) SyncCustomers

```go
func (g *GoConnect) SyncCustomers()
```

#### func (*GoConnect) UserIsAdmin

```go
func (g *GoConnect) UserIsAdmin(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserIsAppUser

```go
func (g *GoConnect) UserIsAppUser(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserIsBot

```go
func (g *GoConnect) UserIsBot(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserIsOwner

```go
func (g *GoConnect) UserIsOwner(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserIsPrimaryOwner

```go
func (g *GoConnect) UserIsPrimaryOwner(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserIsRestricted

```go
func (g *GoConnect) UserIsRestricted(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserIsStranger

```go
func (g *GoConnect) UserIsStranger(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserIsUltraRestricted

```go
func (g *GoConnect) UserIsUltraRestricted(ctx context.Context, email string) (bool, error)
```

#### func (*GoConnect) UserPhoneNumber

```go
func (g *GoConnect) UserPhoneNumber(ctx context.Context, email string) (string, error)
```

#### func (*GoConnect) Util

```go
func (g *GoConnect) Util() *objectify.Handler
```
Util returns an objectify util tool ref:github.com/autom8ter/objectify

#### type LogConfig

```go
type LogConfig struct {
	UserName string `validate:"required"`
	Channel  string `validate:"required"`
}
```


#### type Plugin

```go
type Plugin interface {
	driver.Plugin
	RegisterWithClient(g *GoConnect)
}
```
