# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type Config

```go
type Config struct {
	FirebaseCredsPath string `json:"firebase_creds_path"`
	TwilioAccount     string `json:"twilio_account"`
	TwilioToken       string `json:"twilio_token"`
	SendGridToken     string `json:"sendgrid_token"`
	StripeToken       string `json:"stripe_token"`
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
func New(cli *http.Client, c *Config) *GoConnect
```
New Creates a new GoConnect from the provided http client and config

#### func (*GoConnect) Auth

```go
func (g *GoConnect) Auth() *auth.Client
```
Auth returns an authenticated Firebase Auth client

#### func (*GoConnect) Database

```go
func (g *GoConnect) Database(url string) *db.Client
```
Database returns an authenticated Firebase Database client

#### func (*GoConnect) Firestore

```go
func (g *GoConnect) Firestore() *firestore.Client
```
Firestore returns an authenticated Firebase Firestore client

#### func (*GoConnect) HTTP

```go
func (g *GoConnect) HTTP() *http.Client
```
Twilio returns an HTTP client

#### func (*GoConnect) Messaging

```go
func (g *GoConnect) Messaging() *messaging.Client
```
Messaging returns an authenticated Firebase Messaging client

#### func (*GoConnect) SendGrid

```go
func (g *GoConnect) SendGrid() *sendgrid.Client
```
Twilio returns an authenticated SendGrid client

#### func (*GoConnect) Storage

```go
func (g *GoConnect) Storage() *storage.Client
```
Store returns an authenticated Firebase Storage client

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
