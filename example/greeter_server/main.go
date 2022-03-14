/*
*
* Copyright 2015 gRPC authors.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/chzyer/readline"
	"github.com/lingliy/gpush"
	pb "github.com/lingliy/gpush/example/helloworld"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
	msg gpush.Server
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	msg := gpush.InitServer(s)
	msg.RegisterOnline(func(id string, data []byte) error {
		log.Printf("online %s %s\n", id, data)
		return nil
	})
	msg.RegisterOffline(func(id string) {
		log.Printf("offline %s \n", id)
	})

	srv := &server{}
	pb.RegisterGreeterServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	go func() {
		for {
			line, err := rl.Readline()
			if err != nil { // io.EOF
				break
			}
			args := strings.SplitN(line, " ", 3)
			if len(args) != 3 {
				continue
			}
			log.Println(args)
			msg.SendCmd(args[0], args[1], []byte(args[2]))
		}
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
