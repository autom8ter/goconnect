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
	StripeAccount   string
	StripeToken     string `validate:"required"`
	SlackAccount    string
	SlackToken      string `validate:"required"`
}
```

GoConnect holds the required configuration variables to create a GoConnect
Instance. Use Init() to validate a GoConnect instance.

#### func  NewGoConnect

```go
func NewGoConnect(g *gcloud.GCP, twilioAccount string, twilioToken string, sendGridToken string, stripeAccount string, stripeToken string, slackToken string) *GoConnect
```

#### func (*GoConnect) Execute

```go
func (g *GoConnect) Execute(fns ...HandlerFunc) error
```
Execute runs the provided functions.

#### func (*GoConnect) Gcloud

```go
func (g *GoConnect) Gcloud() *gcloud.GCP
```

#### func (*GoConnect) Init

```go
func (g *GoConnect) Init() error
```

#### func (*GoConnect) SendGrid

```go
func (g *GoConnect) SendGrid() *sendgrid.Client
```

#### func (*GoConnect) Stripe

```go
func (g *GoConnect) Stripe(client *http.Client) *cli.API
```

#### func (*GoConnect) Twilio

```go
func (g *GoConnect) Twilio() *gotwilio.Twilio
```

#### type HandlerFunc

```go
type HandlerFunc func(g *GoConnect) error
```

A HandlerFuncFunc is a GoConnect Callback function handler
