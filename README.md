# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type Config

```go
type Config struct {
	TwilioAccount string `validate:"required"`
	TwilioKey     string `validate:"required"`
	SendgridKey   string `validate:"required"`
	SlackKey      string `validate:"required"`
	StripeKey     string `validate:"required"`
	Index         CustomerIndex
	SyncFrequency string       `validate:"required"`
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
func NewFromEnv(customerIndex CustomerIndex) *GoConnect
```
NewFromFileEnv Initializes a gcp instance from service account credentials ref:
https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
and looks for Twilio, SendGrid, Stripe, and Slack credentials in environmental
variables. vars: TWILIO_ACCOUNT, TWILIO_KEY, SENDGRID_KEY, SYNC_FREQUENCY,
STRIPE_KEY, STRIPE_TOKEN, SLACK_KEY

#### func (*GoConnect) CancelSubscription

```go
func (g *GoConnect) CancelSubscription(key string) error
```

#### func (*GoConnect) Config

```go
func (g *GoConnect) Config() *Config
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

#### func (*GoConnect) Customers

```go
func (g *GoConnect) Customers() map[string]*stripe.Customer
```
Customers returns your current stripe customer cache

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

#### func (*GoConnect) GetSubscription

```go
func (g *GoConnect) GetSubscription(key string) (*stripe.Subscription, bool)
```

#### func (*GoConnect) GetVideoRoom

```go
func (g *GoConnect) GetVideoRoom(id string) (*gotwilio.VideoResponse, error)
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

#### func (*GoConnect) SubscribeCustomer

```go
func (g *GoConnect) SubscribeCustomer(key string, plan, cardnum, month, year, cvc string) (*stripe.Subscription, error)
```

#### func (*GoConnect) Subscriptions

```go
func (g *GoConnect) Subscriptions() map[string]*stripe.Subscription
```

#### func (*GoConnect) SwitchIndex

```go
func (g *GoConnect) SwitchIndex(typ CustomerIndex)
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
