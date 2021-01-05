package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/grpc/greetpb"

	"google.golang.org/grpc"
)

func main() {
	// boilerplate
	fmt.Println("hello I'm a grpc client!")

	cc, err := grpc.Dial("app-mf0gv7g-svc-n-mf0gv70-50051.grpc.local:8080", grpc.WithInsecure())
	defer cc.Close()
	if err != nil {
		log.Fatalln("could not connect:", err)
	}
	c := greetpb.NewGreetServiceClient(cc)

	// invoking grpc calls
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Unary RPC!")
	req := &greetpb.GreetingRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "yang",
			LastName:  "hu",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Printf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.GetResult())
}
