package main

import (
	// "log/slog"

	"log"
	"net"

	"github.com/saladin2098/forum_auth/config"
	pb "github.com/saladin2098/forum_auth/genproto"
	"github.com/saladin2098/forum_auth/service"
	"github.com/saladin2098/forum_auth/storage/postgres"
	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.Load()
	liss, err := net.Listen("tcp", cfg.HTTPPort)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, service.NewService(db))

	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
