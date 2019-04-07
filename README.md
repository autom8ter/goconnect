# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type Config

```go
type Config struct {
	GCPProjectID    string   `validate:"required"`
	GCPCredsPath    string   `validate:"required"`
	TwilioAccount   string   `validate:"required"`
	TwilioToken     string   `validate:"required"`
	SendGridAccount string   `validate:"required"`
	SendGridToken   string   `validate:"required"`
	StripeAccount   string   `validate:"required"`
	StripeToken     string   `validate:"required"`
	SlackAccount    string   `validate:"required"`
	SlackToken      string   `validate:"required"`
	Scopes          []string `validate:"required"`
	InCluster       bool
	MasterKey       string `validate:"required"`
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

#### func (*GoConnect) CompareHashToPassword

```go
func (g *GoConnect) CompareHashToPassword(hash, pass string) error
```
CompareHashToPassWord uses bcrypt to compare the provided hash to the provided
password

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

#### func (*GoConnect) HashPassword

```go
func (g *GoConnect) HashPassword(pass string) (string, error)
```
HashPassword uses bcrypt to hash a password string

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

#### func (*GoConnect) MustGetEnv

```go
func (g *GoConnect) MustGetEnv(key string, defval string) string
```
MustGetEnv returns the environmental variable found in the provided key. If no
value is found, the provided default value is returned

#### func (*GoConnect) NewToken

```go
func (g *GoConnect) NewToken(claims *MyCustomClaims) (string, error)
```
NewToken create a new JWT token from the provided claims

#### func (*GoConnect) Render

```go
func (g *GoConnect) Render(text string, w io.Writer) error
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

#### func (*GoConnect) UUID

```go
func (g *GoConnect) UUID() string
```
Provides UUID string. UUIDs are based on RFC 4122 and DCE 1.1: Authentication
and Security Services.

#### func (*GoConnect) Validate

```go
func (g *GoConnect) Validate(obj interface{}) error
```
Validate validates the provided object and returns an error if invalid ref:
https://github.com/go-playground/validator

#### func (*GoConnect) ValidateToken

```go
func (g *GoConnect) ValidateToken(myToken string) (bool, string)
```
ValidateToken will validate the token

#### func (*GoConnect) WrapErr

```go
func (g *GoConnect) WrapErr(err error, msg string) error
```
WrapErr wraps the provided error with the provided message

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

#### type MyCustomClaims

```go
type MyCustomClaims struct {
	Account string `json:"account"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	jwt.StandardClaims
}
```
