package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/divtxt/lockd/lockapi"
	"github.com/divtxt/lockd/lockimpl"
)

func main() {
	// Parse command line args
	listenAddrPtr := flag.String("listen", ":2080", "listen address")
	flag.Parse()

	// Reset standard log flags (undo Gin's settings)
	log.SetFlags(log.LstdFlags)

	// Instantiate a lock service
	var ilApi lockimpl.InternalLockApi
	var err error
	ilApi, err = lockimpl.NewLockApiImpl()
	if err != nil {
		panic(err)
	}
	lockingServer := lockimpl.NewLockingServerImpl(ilApi)

	// Setup network service
	lis, err := net.Listen("tcp", *listenAddrPtr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	lockapi.RegisterLockingServer(s, lockingServer)

	// Run forever
	log.Println("Starting server on address:", *listenAddrPtr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
