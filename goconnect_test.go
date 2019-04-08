package goconnect_test

import (
	"context"
	"fmt"
	"github.com/autom8ter/goconnect"
	"github.com/autom8ter/goconnect/plugins"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	go goconnect.NewFromFileEnv("credentials.json").Serve(":3000", plugins.EchoService())
	resp, err := dial(context.Background(), "Coleman", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	if resp != "Hola coleman" {
		t.Fatal("incorrect response")
	}
	fmt.Println(resp)
}

func dial(ctx context.Context, name, address string) (string, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	newCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	r, err := c.SayHello(newCtx, &pb.HelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	return r.Message, nil
}