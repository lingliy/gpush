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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/lingliy/gpush"
	pb "github.com/lingliy/gpush/example/helloworld"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	m := gpush.InitClient(conn)
	m.RegisterCmd("/SayHello", func(jsonStr []byte) {
		log.Printf("%s\n", jsonStr)
		c.SayHello(context.Background(), &pb.HelloRequest{Name: string(jsonStr)})
	})
replay:
	err = m.Run("12345")
	if err != nil {
		log.Printf("could not run message: %v\n", err)
		goto replay
	}
}
