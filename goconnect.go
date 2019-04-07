package goconnect

import (
	"context"
	"github.com/autom8ter/gcloud"
	"github.com/autom8ter/goconnect/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"google.golang.org/api/option"
	"io"
	"net/http"
)

//Config holds the required configuration variables to create a GoConnect Instance
type Config struct {
	GCPCredsPath    string   `json:"firebase_creds_path"`
	TwilioAccount   string   `json:"twilio_account"`
	TwilioToken     string   `json:"twilio_token"`
	SendGridAccount string   `json:"sendgrid_account"`
	SendGridToken   string   `json:"sendgrid_token"`
	StripeAccount   string   `json:"stripe_account"`
	StripeToken     string   `json:"stripe_token"`
	SlackAccount    string   `json:"slack_account"`
	SlackToken      string   `json:"slack_token"`
	Scopes          []string `json:"scopes"`
	MasterKey       string   `json:"master_key"`
}

// GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client. It also carries an HTTP client and context.
type GoConnect struct {
	ctx   context.Context
	cfg   *Config
	twil  *gotwilio.Twilio
	grid  *sendgrid.Client
	strip *client.API
	chat  *slack.Client
	gcp   *gcloud.GCP
	cli   *http.Client
	data  map[string]interface{}
}

type MyCustomClaims struct {
	Account string `json:"account"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	jwt.StandardClaims
}

// New Creates a new GoConnect from the provided http client and config
func New(ctx context.Context, c *Config) (*GoConnect, error) {
	g := &GoConnect{}
	var err error
	g.ctx = ctx
	g.cfg = c
	g.gcp, err = gcloud.New(ctx, option.WithCredentialsFile(c.GCPCredsPath))
	if err != nil {
		err = util.WrapErr(err, "failed to create gcloud client")
	}
	g.cli, err = g.gcp.DefaultClient(ctx, c.Scopes)
	if err != nil {
		err = util.WrapErr(err, "failed to create http client from config scopes")
	}
	g.chat = slack.New(c.SlackToken)
	g.twil = gotwilio.NewTwilioClientCustomHTTP(c.TwilioAccount, c.TwilioToken, g.cli)
	g.strip = client.New(c.StripeToken, stripe.NewBackends(g.cli))
	g.grid = sendgrid.NewSendClient(c.SendGridToken)
	g.data = make(map[string]interface{})
	if err != nil {
		return g, err
	}
	return g, nil
}

// GCP returns a gcloud.GCP instance
func (g *GoConnect) GCP() *gcloud.GCP {
	return g.gcp
}

// Config returns the config used to create the GoConnect instance
func (g *GoConnect) Config() *Config {
	return g.cfg
}

// Stripe returns an authenticated Stripe client
func (g *GoConnect) Stripe() *client.API {
	return g.strip
}

// Twilio returns an authenticated Twilio client
func (g *GoConnect) Twilio() *gotwilio.Twilio {
	return g.twil
}

// SendGrid returns an authenticated SendGrid client
func (g *GoConnect) SendGrid() *sendgrid.Client {
	return g.grid
}

// Slack returns an authenticated Slack client
func (g *GoConnect) Slack() *slack.Client {
	return g.chat
}

// HTTP returns an HTTP client
func (g *GoConnect) HTTP() *http.Client {
	return g.cli
}

// Render renders the text with the GoConnects current data. It writes the output to the provided writer
func (g *GoConnect) Render(text string, w io.Writer) error {
	return util.Render(text, g.data, w)
}

// JSON returns the GoConnects current data as JSON
func (g *GoConnect) JSON() []byte {
	return util.JSON(g.data)
}

// YAML returns the GoConnects current data as YAML
func (g *GoConnect) YAML() []byte {
	return util.YAML(g.data)
}

// XML returns the GoConnects current data as XML
func (g *GoConnect) XML() []byte {
	return util.XML(g.data)
}

// Validate validates the provided object and returns an error if invalid ref: https://github.com/go-playground/validator
func (g *GoConnect) Validate(obj interface{}) error {
	return util.Validate(obj)
}

// Data returns GoConnects current data as map[string]interface{}
func (g *GoConnect) Data() map[string]interface{} {
	return g.data
}

// AddStructData appends the provided structs data to the GoConnects data
func (g *GoConnect) AddStructData(obj interface{}) {
	for k, v := range util.AsMap(obj) {
		g.data[k] = v
	}
}

// AddData appends the provided data to the GoConnects current data
func (g *GoConnect) AddData(obj map[string]interface{}) {
	for k, v := range obj {
		g.data[k] = v
	}
}

// MasterKey returns the master key from config as bytes. Defaults to "secret"
func (g *GoConnect) MasterKey() []byte {
	if g.cfg.MasterKey != "" {
		return []byte(g.cfg.MasterKey)
	}
	return []byte("secret")
}

// HashPassword uses bcrypt to hash a password string
func (g *GoConnect) HashPassword(pass string) (string, error) {
	return util.HashPassword(pass)
}

// CompareHashToPassWord uses bcrypt to compare the provided hash to the provided password
func (g *GoConnect) CompareHashToPassword(hash, pass string) error {
	return util.ComparePasswordToHash(hash, pass)
}

// WrapErr wraps the provided error with the provided message
func (g *GoConnect) WrapErr(err error, msg string) error {
	return util.WrapErr(err, msg)
}

// Provides UUID string. UUIDs are based on RFC 4122 and DCE 1.1: Authentication and Security Services.
func (g *GoConnect) UUID() string {
	return util.UUID()
}

// MustGetEnv returns the environmental variable found in the provided key. If no value is found, the provided default value is returned
func (g *GoConnect) MustGetEnv(key string, defval string) string {
	return util.MustGetEnv(key, defval)
}

// NewToken create a new JWT token from the provided claims
func (g *GoConnect) NewToken(claims *MyCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	/* Sign the token with our secret */
	return token.SignedString(g.MasterKey())
}

// ValidateToken will validate the token
func (g *GoConnect) ValidateToken(myToken string) (bool, string) {
	token, err := jwt.ParseWithClaims(myToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(g.MasterKey()), nil
	})

	if err != nil {
		return false, ""
	}

	claims := token.Claims.(*MyCustomClaims)
	return token.Valid, claims.Email
}

// A Func is a GoConnect Callback function
type Func func(g *GoConnect) error

// Execute runs the provided functions.
func (g *GoConnect) Execute(fns ...Func) error {
	var err error
	for _, f := range fns {
		if newErr := f(g); newErr != nil {
			err = errors.Wrap(err, newErr.Error())
		}
	}
	return err
}
