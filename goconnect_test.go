package goconnect_test

import (
	"github.com/autom8ter/goconnect"
	"log"
	"os"
	"testing"
)

func init() {
	g = goconnect.New(nil, &goconnect.Config{
		GCPCredsPath:  "credentials.json",
		TwilioAccount: os.Getenv("TWILIO_ACCOUNT"),
		TwilioToken:   os.Getenv("TWILIO_TOKEN"),
		SendGridToken: os.Getenv("SENDGRID_TOKEN"),
		StripeToken:   os.Getenv("STRIPE_TOKEN"),
	})
}

var g *goconnect.GoConnect

func TestNewFromEnv(t *testing.T) {
	if g == nil {
		log.Fatalln("nil goconnect")
	}
}
