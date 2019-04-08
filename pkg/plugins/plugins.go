package plugins

import (
	"context"
	"github.com/autom8ter/engine/driver"
	"github.com/autom8ter/goconnect"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct {
	g *goconnect.GoConnect
}

func NewServer(g *goconnect.GoConnect) *Server {
	return &Server{g: g}
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	cli, err := s.g.GCP.Translate(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := cli.Translate(ctx, []string{"Hello " + in.Name}, language.Spanish, nil)

	return &pb.HelloReply{Message: resp[0].Text}, nil
}

func GreeterService() goconnect.PluginFunc {
	return func(g *goconnect.GoConnect) driver.PluginFunc {
		return func(s *grpc.Server) {
			pb.RegisterGreeterServer(s, &Server{
				g,
			})
		}
	}
}

func TranslationService() goconnect.PluginFunc {
	return func(g *goconnect.GoConnect) driver.PluginFunc {
		return func(s *grpc.Server) {
			pb.RegisterGreeterServer(s, &Server{
				g,
			})
		}
	}
}
