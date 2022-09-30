package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen to port  9000 %v", err)
	}

	grcpServer := grpc.NewServer()
	if err := grcpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve  grcpServer over port 9000 %v", err)
	}
}
