package plugins

import (
	"context"
	"github.com/autom8ter/engine/driver"
	"github.com/autom8ter/goconnect"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func EchoService() goconnect.PluginFunc {
	return func(g *goconnect.GoConnect) driver.PluginFunc {
		return func(s *grpc.Server) {
			pb.RegisterGreeterServer(s, &server{})
		}
	}
}
