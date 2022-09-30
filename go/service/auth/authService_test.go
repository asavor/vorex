package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	err := godotenv.Load("env")
	if err != nil {
		log.Fatalln("Failed to load .ENV file.")
		return
	}
	ctx := context.Background()

	rDB := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	user, err := rDB.Get(ctx, "user1").Result()

	if err != nil {
		fmt.Println("can not find user")
		return
	}
	fmt.Println("my name is " + user)

}
