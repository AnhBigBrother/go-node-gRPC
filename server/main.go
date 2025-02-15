package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	pb "greating-grpc/proto"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received from client : %s", in.GetMessage())
	fmt.Println("-----")
	return &pb.HelloReply{Message: "Hello from server!"}, nil
}

func (s *server) SayHelloStreamRequest(stream pb.Greeter_SayHelloStreamRequestServer) error {
	count := 0
	msgs := []string{}
	for {
		greating, err := stream.Recv()
		if err == io.EOF {
			msg := fmt.Sprintf("Received %d messages from client: %s", count, strings.Join(msgs, ", "))
			fmt.Println(msg)
			fmt.Println("-----")
			return stream.SendAndClose(&pb.HelloReply{
				Message: msg,
			})
		}
		if err != nil {
			return err
		}
		msgs = append(msgs, greating.GetMessage())
		fmt.Println(greating)
		count++
	}
}

func (s *server) SayHelloStreamReply(in *pb.HelloRequest, stream pb.Greeter_SayHelloStreamReplyServer) error {
	fmt.Println("Received from client:", in.GetMessage())
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Hello from server no.%d", i)
		if err := stream.Send(&pb.HelloReply{Message: msg}); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Say hello to client 10 times")
	fmt.Println("-----")
	return nil
}

func (s *server) SayHelloBidirectionalStreaming(stream pb.Greeter_SayHelloBidirectionalStreamingServer) error {
	count := 1
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("-----")
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println("Received from client:", in.Message)
		stream.Send(&pb.HelloReply{
			Message: fmt.Sprintf("Server received %d messages", count),
		})
		count++
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
