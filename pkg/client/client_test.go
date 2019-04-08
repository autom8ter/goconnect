package client_test

import (
	"context"
	"github.com/autom8ter/gosaas/pkg/client"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"google.golang.org/grpc"
	"os"
	"testing"
)

var ctx = context.Background()

func Test(t *testing.T) {
	c := client.New(ctx, "localhost:3000", grpc.WithInsecure())
	_, err := c.SMS().SendSMS(ctx, &contacts.SMS{
		From: &contacts.Contact{
			Phone: os.Getenv("TWILIO_PHONE"),
		},
		To: &contacts.Contact{
			Phone: os.Getenv("TWILIO_TEST"),
		},
		Body: "testing GoSaaS client",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}
