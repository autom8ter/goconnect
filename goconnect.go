package goconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/autom8ter/gcloud"
	"github.com/autom8ter/objectify"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
)

var tool = objectify.New()

//Config holds the required configuration variables to create a GoConnect Instance
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

// GoConnect holds an authenticated Twilio, Stripe, Firebase, and SendGrid Client. It also carries an HTTP client and context.
type GoConnect struct {
	ctx   context.Context        `validate:"required"`
	cfg   *Config                `validate:"required"`
	twil  *gotwilio.Twilio       `validate:"required"`
	grid  *sendgrid.Client       `validate:"required"`
	strip *client.API            `validate:"required"`
	chat  *slack.Client          `validate:"required"`
	gcp   *gcloud.GCP            `validate:"required"`
	data  map[string]interface{} `validate:"required"`
}

// New Creates a new GoConnect from the provided http client and config
func New(ctx context.Context, c *Config) (*GoConnect, error) {
	if err := tool.Validate(c); err != nil {
		panic(err.Error())
	}
	gcp := gcloud.New(ctx, &gcloud.Config{
		Project:   c.ProjectID,
		Scopes:    c.Scopes,
		InCluster: false,
		Options:   []option.ClientOption{option.WithCredentialsFile(c.JSONPath)},
	})
	if err := gcp.Error(); err != nil {
		log.Println("gcp initialization error- ignore if not using client in error.", err.Error())
	}

	chat := slack.New(c.SlackToken)
	twil := gotwilio.NewTwilioClientCustomHTTP(c.TwilioAccount, c.TwilioToken, gcp.HTTP())
	strip := client.New(c.StripeToken, stripe.NewBackends(gcp.HTTP()))
	grid := sendgrid.NewSendClient(c.SendGridToken)
	data := make(map[string]interface{})
	g := &GoConnect{
		ctx:   ctx,
		cfg:   c,
		twil:  twil,
		grid:  grid,
		chat:  chat,
		strip: strip,
		data:  data,
	}
	if err := tool.Validate(g); err != nil {
		panic(err.Error())
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
	return g.gcp.HTTP()
}

// Render renders the text with the GoConnects current data. It writes the output to the provided writer
func (g *GoConnect) RenderTXT(text string, w io.Writer) error {
	return tool.RenderTXT(text, g.data, w)
}

// Render renders the text with the GoConnects current data. It writes the output to the provided writer
func (g *GoConnect) RenderHTML(text string, w io.Writer) error {
	return tool.RenderHTML(text, g.data, w)
}

// JSON returns the GoConnects current data as JSON
func (g *GoConnect) JSON() []byte {
	return tool.MarshalJSON(g.data)
}

// YAML returns the GoConnects current data as YAML
func (g *GoConnect) YAML() []byte {
	return tool.MarshalYAML(g.data)
}

// XML returns the GoConnects current data as XML
func (g *GoConnect) XML() []byte {
	return tool.MarshalXML(g.data)
}

// Data returns GoConnects current data as map[string]interface{}
func (g *GoConnect) Data() map[string]interface{} {
	return g.data
}

// AddStructData appends the provided structs data to the GoConnects data
func (g *GoConnect) AddStructData(obj interface{}) {
	for k, v := range tool.ToMap(obj) {
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

// A HandlerFuncFunc is a GoConnect Callback function handler
type HandlerFunc func(g *GoConnect) error

// Execute runs the provided functions.
func (g *GoConnect) Execute(fns ...HandlerFunc) error {
	var err error
	for _, f := range fns {
		if newErr := f(g); newErr != nil {
			err = errors.Wrap(err, newErr.Error())
		}
	}
	return err
}

func toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
}
