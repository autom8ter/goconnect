# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type GoConnect

```go
type GoConnect struct {
	GCP             *gcloud.GCP `validate:"required"`
	TwilioAccount   string      `validate:"required"`
	TwilioToken     string      `validate:"required"`
	SendGridAccount string
	SendGridToken   string `validate:"required"`
	StripeAccount   string `validate:"required"`
	StripeToken     string `validate:"required"`
	SlackAccount    string
	SlackToken      string `validate:"required"`
}
```

GoConnect holds the required configuration variables to create a GoConnect
Instance. Use Init() to validate a GoConnect instance.

#### func  New

```go
func New(twilioAccount string, twilioToken string, sendGridToken string, stripeAccount string, stripeToken string, slackToken string, opts ...option.ClientOption) *GoConnect
```
New creates a new GoConnect Instance (no magic)

#### func  NewFromFileEnv

```go
func NewFromFileEnv(file string) *GoConnect
```
NewFromFileEnv Initializes a gcp instance from service account credentials ref:
https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
and looks for Twilio, SendGrid, Stripe, and Slack credentials in environmental
variables. vars: TWILIO_ACCOUNT, TWILIO_ACCOUNT, SENDGRID_ACCOUNT,
SENDGRID_TOKEN, STRIPE_ACCOUNT, STRIPE_TOKEN, SLACK_ACCOUNT, SLACK_TOKEN

#### func (*GoConnect) Gcloud

```go
func (g *GoConnect) Gcloud() *gcloud.GCP
```
Gcloud returns an authenticated GCP instance

#### func (*GoConnect) Init

```go
func (g *GoConnect) Init() error
```
Init returns an error if any of the required fields are nil

#### func (*GoConnect) SendGrid

```go
func (g *GoConnect) SendGrid() *sendgrid.Client
```
SendGrid returns an authenticated SendGrid client

#### func (*GoConnect) Serve

```go
func (g *GoConnect) Serve(addr string, fns ...PluginFunc) error
```
Serve starts a grpc Engine ref:github.com/autom8ter/engine server with a default
middleware stack on the specified address with the provided pluugin functions.

#### func (*GoConnect) Stripe

```go
func (g *GoConnect) Stripe(client *http.Client) *cli.API
```
Stripe returns an authenticated Stripe client

#### func (*GoConnect) ToMap

```go
func (g *GoConnect) ToMap() map[string]interface{}
```
ToMap returns the GoConnect config as a map

#### func (*GoConnect) Twilio

```go
func (g *GoConnect) Twilio() *gotwilio.Twilio
```
Twilio returns an authenticated Twilio client

#### type PluginFunc

```go
type PluginFunc func(g *GoConnect) driver.PluginFunc
```

PluginFunc is a callback function that takes a GoConnect instance and returns a
function that is used to create and register a grpc service. It is used in the
GoConnect Serve() method.
