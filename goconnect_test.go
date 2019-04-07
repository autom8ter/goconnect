package goconnect_test

import (
	"context"
	"github.com/autom8ter/goconnect"
	"log"
	"os"
	"testing"
)

var ctx  = context.Background()
var err error
var g *goconnect.GoConnect

func init() {
	g, err = goconnect.New(ctx,  &goconnect.Config{
		GCPCredsPath:  "credentials.json",
		TwilioAccount: os.Getenv("TWILIO_ACCOUNT"),
		TwilioToken:   os.Getenv("TWILIO_TOKEN"),
		SendGridToken: os.Getenv("SENDGRID_TOKEN"),
		StripeToken:   os.Getenv("STRIPE_TOKEN"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}


func TestNewFromEnv(t *testing.T) {
	if g == nil {
		log.Fatalln("nil goconnect")
	}
}
