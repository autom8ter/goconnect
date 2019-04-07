# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type Config

```go
type Config struct {
	ProjectID       string `validate:"required"`
	JSONPath        string `validate:"required"`
	TwilioAccount   string `validate:"required"`
	TwilioToken     string `validate:"required"`
	SendGridAccount string
	SendGridToken   string `validate:"required"`
	StripeAccount   string
	StripeToken     string `validate:"required"`
	SlackAccount    string
	SlackToken      string   `validate:"required"`
	Scopes          []string `validate:"required"`
	InCluster       bool
	MasterKey       string
}
```

Config holds the required configuration variables to create a GoConnect Instance

#### type GoConnect

```go
type GoConnect struct {
}
```

GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client.
It also carries an HTTP client and context.

#### func  New

```go
func New(ctx context.Context, c *Config) (*GoConnect, error)
```
New Creates a new GoConnect from the provided http client and config

#### func (*GoConnect) AddData

```go
func (g *GoConnect) AddData(obj map[string]interface{})
```
AddData appends the provided data to the GoConnects current data

#### func (*GoConnect) AddStructData

```go
func (g *GoConnect) AddStructData(obj interface{})
```
AddStructData appends the provided structs data to the GoConnects data

#### func (*GoConnect) Config

```go
func (g *GoConnect) Config() *Config
```
Config returns the config used to create the GoConnect instance

#### func (*GoConnect) Data

```go
func (g *GoConnect) Data() map[string]interface{}
```
Data returns GoConnects current data as map[string]interface{}

#### func (*GoConnect) Execute

```go
func (g *GoConnect) Execute(fns ...HandlerFunc) error
```
Execute runs the provided functions.

#### func (*GoConnect) GCP

```go
func (g *GoConnect) GCP() *gcloud.GCP
```
GCP returns a gcloud.GCP instance

#### func (*GoConnect) HTTP

```go
func (g *GoConnect) HTTP() *http.Client
```
HTTP returns an HTTP client

#### func (*GoConnect) JSON

```go
func (g *GoConnect) JSON() []byte
```
JSON returns the GoConnects current data as JSON

#### func (*GoConnect) MasterKey

```go
func (g *GoConnect) MasterKey() []byte
```
MasterKey returns the master key from config as bytes. Defaults to "secret"

#### func (*GoConnect) RenderHTML

```go
func (g *GoConnect) RenderHTML(text string, w io.Writer) error
```
Render renders the text with the GoConnects current data. It writes the output
to the provided writer

#### func (*GoConnect) RenderTXT

```go
func (g *GoConnect) RenderTXT(text string, w io.Writer) error
```
Render renders the text with the GoConnects current data. It writes the output
to the provided writer

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

#### func (*GoConnect) XML

```go
func (g *GoConnect) XML() []byte
```
XML returns the GoConnects current data as XML

#### func (*GoConnect) YAML

```go
func (g *GoConnect) YAML() []byte
```
YAML returns the GoConnects current data as YAML

#### type HandlerFunc

```go
type HandlerFunc func(g *GoConnect) error
```

A HandlerFuncFunc is a GoConnect Callback function handler
