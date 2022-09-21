package main

import (
	"context"
	"flag"
	"github.com/openziti-test-kitchen/grpc-ziti-starter/protocol"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"google.golang.org/grpc"
	"log"
)

var (
	identity = flag.String("identity", "", "Ziti Identity file")
	service  = flag.String("service", "", "Ziti Service")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	protocol.UnimplementedAnswerServiceServer
}

// WhatIs implements protocol.AnswerServiceServer
func (s *server) WhatIs(_ context.Context, in *protocol.Question) (*protocol.Answer, error) {
	term := in.GetWhat()

	log.Printf("Question: what is %v?", term)
	if term == "ziti" {
		return &protocol.Answer{Answer: "ziti is a type of pasta"}, nil
	}
	return &protocol.Answer{Answer: "I don't know what " + in.GetWhat() + " is :("}, nil
}

func main() {
	flag.Parse()
	cfg, err := config.NewFromFile(*identity)
	if err != nil {
		log.Fatalf("failed to load ziti identity{%v}: %v", identity, err)
	}

	ztx := ziti.NewContextWithConfig(cfg)
	err = ztx.Authenticate()
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}

	lis, err := ztx.Listen(*service)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protocol.RegisterAnswerServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
