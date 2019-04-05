# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type Config

```go
type Config struct {
	GCPCredsPath    string `json:"firebase_creds_path"`
	TwilioAccount   string `json:"twilio_account"`
	TwilioToken     string `json:"twilio_token"`
	SendGridAccount string `json:"sendgrid_account"`
	SendGridToken   string `json:"sendgrid_token"`
	StripeAccount   string `json:"stripe_account"`
	StripeToken     string `json:"stripe_token"`
	SlackAccount    string `json:"slack_account"`
	SlackToken      string `json:"slack_token"`
}
```

Config holds the required configuration variables to create a GoConnect Instance

#### type Func

```go
type Func func(g *GoConnect) error
```


#### type GoConnect

```go
type GoConnect struct {
}
```

GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client.
It also carries an HTTP client and context.

#### func  New

```go
func New(cli *http.Client, c *Config) *GoConnect
```
New Creates a new GoConnect from the provided http client and config

#### func (*GoConnect) Config

```go
func (g *GoConnect) Config() *Config
```
Config returns the config used to create the GoConnect instance

#### func (*GoConnect) Execute

```go
func (g *GoConnect) Execute(fns ...Func) error
```

#### func (*GoConnect) GCP

```go
func (g *GoConnect) GCP() *gcloud.GCP
```
GCP returns a gcloud.GCP instance

#### func (*GoConnect) GoCrypt

```go
func (g *GoConnect) GoCrypt() *gocrypt.GoCrypt
```
GoSub returns a GoCrypt client

#### func (*GoConnect) HTTP

```go
func (g *GoConnect) HTTP() *http.Client
```
Twilio returns an HTTP client

#### func (*GoConnect) SendGrid

```go
func (g *GoConnect) SendGrid() *sendgrid.Client
```
SendGrid returns an authenticated SendGrid client

#### func (*GoConnect) Slack

```go
func (g *GoConnect) Slack() *slack.Client
```
Slack returns an authenticated Slack client

#### func (*GoConnect) Stripe

```go
func (g *GoConnect) Stripe() *client.API
```
Stripe returns an authenticated Stripe client

#### func (*GoConnect) Twilio

```go
func (g *GoConnect) Twilio() *gotwilio.Twilio
```
Twilio returns an authenticated Twilio client
