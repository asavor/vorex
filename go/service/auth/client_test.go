package main

import (
	pb "auth/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return
	}

	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.Login(ctx, &pb.LoginRequest{
		Email:     "helllo@outlook.com",
		Password:  "dsadsad",
		ReCaptcha: "recapToken",
	})

	fmt.Println(resp)

}

func TestSignUp(t *testing.T) {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return
	}

	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.Register(ctx, &pb.RegisterRequest{
		Email:     "helllo@outlook.com",
		Password:  "dsadsad",
		ReCaptcha: "recapToken",
	})

	fmt.Println(resp)

}
