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

#### type CustomerInfo

```go
type CustomerInfo struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	Plans       []*Plan           `json:"plan"`
	Annotations map[string]string `json:"annotations"`
}
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

#### func (*GoConnect) CallCustomer

```go
func (g *GoConnect) CallCustomer(customerKey, from, callback string) (*gotwilio.VoiceResponse, error)
```
Call calls a number

#### func (*GoConnect) CancelSubscription

```go
func (g *GoConnect) CancelSubscription(key string) error
```

#### func (*GoConnect) Compile

```go
func (g *GoConnect) Compile(c *CustomerInfo, hTML string, w io.Writer) error
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

#### func (*GoConnect) Customers

```go
func (g *GoConnect) Customers() map[string]*stripe.Customer
```
Customers returns your current stripe customer cache

#### func (*GoConnect) EmailCustomer

```go
func (g *GoConnect) EmailCustomer(customerKey, subject, plain, html string) error
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

#### func (*GoConnect) ToCustomer

```go
func (g *GoConnect) ToCustomer(c *CustomerInfo) (*stripe.Customer, error)
```

#### func (*GoConnect) ToCustomerInfo

```go
func (g *GoConnect) ToCustomerInfo(customerKey string) (*CustomerInfo, error)
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


#### type Plan

```go
type Plan struct {
	Name    string `json:"id"`
	Amount  int64  `json:"amount"`
	Active  bool   `json:"active"`
	Service string `json:"service"`
}
```


#### type Plugin

```go
type Plugin interface {
	driver.Plugin
	RegisterWithClient(g *GoConnect)
}
```
