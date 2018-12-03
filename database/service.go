package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"../services"
	"google.golang.org/grpc"
)

// DatabaseService is responsibe for storing records in database
type DatabaseService struct {
	db Database
}

// AddRecord calls database to store record
func (s *DatabaseService) AddRecord(ctx context.Context, r *services.Record) (*services.Nothing, error) {
	s.db.AddRecord(*r)
	return &services.Nothing{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	service := DatabaseService{NewMemoryDatabase()}
	services.RegisterDatabaseServiceServer(server, &service)

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
