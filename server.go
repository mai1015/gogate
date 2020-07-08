package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	pb "github.com/mai1015/gogate/rpc"
)

var (
	rpcAddr = flag.String("rpc-addr", ":9080", "RPC TCP address to listen to")
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}


func main() {
	flag.Parse()
	// Create network listener for accepting incoming requests.
	//
	// Note that you are not limited by TCP listener - arbitrary
	// net.Listener may be used by the server.
	// For example, unix socket listener or TLS listener.
	ln, err := net.Listen("tcp4", *addr)
	if err != nil {
		log.Fatalf("error in net.Listen: %s", err)
	}

	// This function will be called by the server for each incoming request.
	//
	// RequestCtx provides a lot of functionality related to http request
	// processing. See RequestCtx docs for details.

	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// Serve returns on ln.Close() or error, so usually it blocks forever.
	r := NewRouting()

	h := r.Handler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}
	h = CombinedColored(h)
	if err := fasthttp.Serve(ln, h); err != nil {
		log.Fatalf("error in Serve: %s", err)
	}
	log.Print("http start at", addr)
	startRPC()
}

func startRPC() {
	lis, err := net.Listen("tcp", *rpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Print("rpc start at", rpcAddr)
}
