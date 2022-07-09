package main

import (
	"back/proto"
	"back/repo"
	"back/server"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	r, err := repo.New()
	if err != nil {
		log.Fatalf("repo error: %s", err)
	}

	err = r.Migrate()
	if err != nil {
		log.Fatalf("migration error: %s", err)
	}

	s := server.NewServer(r)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", os.Getenv("API_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterAPIServer(srv, s)
	log.Printf("server listening at %v", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
