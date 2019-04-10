# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type Config

```go
type Config struct {
	TwilioAccount string        `validate:"required"`
	TwilioKey     string        `validate:"required"`
	SendgridKey   string        `validate:"required"`
	SlackKey      string        `validate:"required"`
	StripeKey     string        `validate:"required"`
	Index         CustomerIndex `validate:"required"`
	SyncFrequency string        `validate:"required"`
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
func NewFromEnv() *GoConnect
```
NewFromFileEnv Initializes a gcp instance from service account credentials ref:
https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
and looks for Twilio, SendGrid, Stripe, and Slack credentials in environmental
variables. vars: TWILIO_ACCOUNT, TWILIO_KEY, SENDGRID_KEY, SYNC_FREQUENCY,
STRIPE_KEY, STRIPE_TOKEN, SLACK_KEY

#### func (*GoConnect) Call

```go
func (g *GoConnect) Call(from, to, callback string) (*gotwilio.VoiceResponse, error)
```
Call calls a number

#### func (*GoConnect) Customers

```go
func (g *GoConnect) Customers() map[string]*stripe.Customer
```
Customers returns your current stripe customer cache

#### func (*GoConnect) GetCustomer

```go
func (g *GoConnect) GetCustomer(key string) *stripe.Customer
```

#### func (*GoConnect) GetSubscription

```go
func (g *GoConnect) GetSubscription(key string) *stripe.Subscription
```

#### func (*GoConnect) Init

```go
func (g *GoConnect) Init() error
```
Init starts syncing the customer cache and validates the GoConnect instance

#### func (*GoConnect) SMS

```go
func (g *GoConnect) SMS(from, to, body, mediaUrl, callback, app string) (*gotwilio.SmsResponse, error)
```
SendSMS sends an sms if mediaurl if empty, mms otherwise.

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
