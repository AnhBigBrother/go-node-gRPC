package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	pb "greating-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received from client : %s", in.GetMessage())
	log.Println("-----")
	return &pb.HelloReply{Message: "Hello from server!"}, nil
}

func (s *server) SayHelloStreamRequest(stream pb.Greeter_SayHelloStreamRequestServer) error {
	count := 0
	msgs := []string{}
	for {
		greating, err := stream.Recv()
		if err == io.EOF {
			msg := fmt.Sprintf("Received %d messages from client: %s", count, strings.Join(msgs, ", "))
			log.Println(msg)
			log.Println("-----")
			return stream.SendAndClose(&pb.HelloReply{
				Message: msg,
			})
		}
		if err != nil {
			return err
		}
		msgs = append(msgs, greating.GetMessage())
		log.Println(greating)
		count++
	}
}

func (s *server) SayHelloStreamReply(in *pb.HelloRequest, stream pb.Greeter_SayHelloStreamReplyServer) error {
	log.Println("Received from client:", in.GetMessage())
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Hello from server no.%d", i)
		if err := stream.Send(&pb.HelloReply{Message: msg}); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	log.Println("Say hello to client 10 times")
	log.Println("-----")
	return nil
}

func (s *server) SayHelloBidirectionalStreaming(stream pb.Greeter_SayHelloBidirectionalStreamingServer) error {
	count := 1
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println("-----")
			return nil
		}
		if err != nil {
			return err
		}

		log.Println("Received from client:", in.Message)
		stream.Send(&pb.HelloReply{
			Message: fmt.Sprintf("Server received %d messages", count),
		})
		count++
	}
}

func unaryServerInterceptorFunc(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Println("-----")
	log.Println("unary interceptor:", info.FullMethod)
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println("authorization", md["authorization"])
	log.Println("description", md["description"])

	return handler(ctx, req)
}

func streamServerInterceptorFunc(
	srv any,
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("-----")
	log.Println("stream interceptor:", info.FullMethod)
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Println("authorization", md["authorization"])
	log.Println("description", md["description"])

	return handler(srv, stream)
}

func loadTSLCredentials() (credentials.TransportCredentials, error) {
	// load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair("../cert/server-cert.pem", "../cert/server-key.pem")
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsCredentials, err := loadTSLCredentials()
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(unaryServerInterceptorFunc),
		grpc.StreamInterceptor(streamServerInterceptorFunc),
	)

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
