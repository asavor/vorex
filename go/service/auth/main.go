package main

import (
	"auth/handler"
	pb "auth/proto"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var ctx = context.Background()

const (
	port = ":50051"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Failed to load .ENV file.")
		return
	}

	rDB := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
	_, err = rDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to redis!")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("fail to listen %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterAuthServiceServer(s, &handler.AuthServiceServer{
		RedisDB: rDB,
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("fail to server %v", err)
	}

}
