# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type CallbackFunc

```go
type CallbackFunc func(customer2 *stripe.Customer) error
```


#### type Config

```go
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
TWILIO_ACCOUNT, TWILIO_KEY, SENDGRID_KEY, SLACK_KEY, STRIPE_KEY, EMAIL_ADDRESS,
EMAIL_NAME, SLACK_LOG_USERNAME, SLACK_LOG_CHANNEL

#### func (*GoConnect) CallBack

```go
func (g *GoConnect) CallBack(key string, funcs ...CallbackFunc) error
```

#### func (*GoConnect) CancelSubscription

```go
func (g *GoConnect) CancelSubscription(key string) error
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

#### func (*GoConnect) CustomerExists

```go
func (g *GoConnect) CustomerExists(key string) bool
```

#### func (*GoConnect) CustomerKeys

```go
func (g *GoConnect) CustomerKeys() []string
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

#### func (*GoConnect) GetSubscriptionFromCustomerEmail

```go
func (g *GoConnect) GetSubscriptionFromCustomerEmail(email string) *stripe.Subscription
```

#### func (*GoConnect) GetSubscriptionFromCustomerID

```go
func (g *GoConnect) GetSubscriptionFromCustomerID(id string) *stripe.Subscription
```

#### func (*GoConnect) GetSubscriptionFromCustomerPhone

```go
func (g *GoConnect) GetSubscriptionFromCustomerPhone(phone string) *stripe.Subscription
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

#### func (*GoConnect) SwitchIndex

```go
func (g *GoConnect) SwitchIndex(typ CustomerIndex)
```

#### func (*GoConnect) SyncCustomers

```go
func (g *GoConnect) SyncCustomers()
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
