package server_test

import (
	"context"
	"github.com/autom8ter/gosaas/pkg/client"
	"github.com/autom8ter/gosaas/pkg/server"
	"github.com/autom8ter/objectify"
	"google.golang.org/grpc"
	"testing"
)

var ctx = context.Background()
var util = objectify.New()
var serv = server.New(ctx, "credentials.json", "autom8ter19", "default")
var cli = client.New(ctx, "localhost:3000", grpc.WithInsecure())

func Test(t *testing.T) {
	t.Fatal(serv.Serve("tcp", ":3000", true))
}

/*
_, err := cli.SMS().SendSMS(ctx, &contacts.SMS{
		From: &contacts.Contact{
			Phone: "+17209032013",
		},
		To: &contacts.Contact{
			Phone: "+13038752836",
		},
		Body: "testing GoSaas. " + time.Now().String(),
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	resp, err := cli.Accounts().ReadAccount(ctx, &accounts.ReadAccountRequest{
		Email: "colemanword@gmail.com",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(string(util.MarshalJSON(resp)))
*/
