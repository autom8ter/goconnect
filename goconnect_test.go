package goconnect_test

import (
	"context"
	"github.com/autom8ter/goconnect"
	"log"
	"os"
	"testing"
)

var ctx = context.Background()
var err error
var g *goconnect.GoConnect

func init() {
	g, err = goconnect.New(ctx, &goconnect.Config{
		ProjectID:       os.Getenv("PROJECT_ID"),
		JSONPath:        "credentials.json",
		TwilioAccount:   os.Getenv("TWILIO_ACCOUNT"),
		TwilioToken:     os.Getenv("TWILIO_TOKEN"),
		SendGridAccount: os.Getenv("SENDGRID_ACCOUNT"),
		SendGridToken:   os.Getenv("SENDGRID_TOKEN"),
		StripeAccount:   os.Getenv("STRIPE_ACCOUNT"),
		StripeToken:     os.Getenv("STRIPE_TOKEN"),
		SlackAccount:    os.Getenv("SLACK_ACCOUNT"),
		SlackToken:      os.Getenv("SLACK_TOKEN"),
		Scopes:          []string{"users"},
		InCluster:       false,
		MasterKey:       os.Getenv("PROJECT_ID"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestNewFromEnv(t *testing.T) {
	if g == nil {
		t.Fatal("nil goconnect")
	}
}
