package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"

	"fmt"

	"github.com/divtxt/lockd/ginx"
	"github.com/divtxt/lockd/httpimpl"
	"github.com/divtxt/lockd/lockimpl"
	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft"
	"github.com/gin-gonic/gin"
)

func main() {
	// Reset standard log flags (undo Gin's settings)
	log.SetFlags(log.LstdFlags)

	// Disable Gin debug logging
	gin.SetMode(gin.ReleaseMode)

	// Parse command line args
	args, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		flag.Usage()
		os.Exit(2)
	}

	// Process args
	log.Println("lockd initializing")
	log.Println("thisServerId:", args.thisServerId)
	cd, err := util.LoadClusterDefinition(args.cluster)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading cluster definition:", err)
		os.Exit(2)
	}
	log.Println("cluster:", cd)

	// Calculate lisen host port using port from cluster info
	_, listenPort, err := net.SplitHostPort(cd.GetHostPort(args.thisServerId))
	listenAddr := fmt.Sprintf("localhost:%s", listenPort)
	log.Println("calculated listenAddr:", listenAddr)

	// Instantiate a lock service
	var l httpimpl.LockApi
	var rcm raft.IConsensusModule
	l, rcm, err = lockimpl.NewLockApiImpl(cd, args.thisServerId)
	if err != nil {
		panic(err)
	}

	// Configure http service
	r := gin.New()
	r.Use(ginx.StdLogLogger())
	r.Use(ginx.StdLogRepanic())

	httpimpl.AddRaftRpcEndpoints(r, rcm)
	httpimpl.AddLockApiEndpoints(r, l)

	// Run forever / till stopped
	log.Println("Starting server on address:", listenAddr)
	err = r.Run(listenAddr)
	if err != nil {
		panic(err)
	}
}

type lockdArgs struct {
	cluster      string
	thisServerId raft.ServerId
}

func parseArgs() (*lockdArgs, error) {
	var args lockdArgs

	// Parse command line args
	flag.StringVar(
		&args.cluster,
		"cluster",
		"",
		"cluster definition file - json file that describes the lockd cluster\n"+
			"    \t    the json should be of the form: {\"server-id\": \"host:port\", ...}\n"+
			"    \t    server ids should be positive integers, but as strings since json keys must be strings\n"+
			"    \t    example: {\"1\": \"lockd1:2081\", \"2\": \"lockd2:2082\", \"3\": \"lockd3:2083\"}",
	)
	flag.Uint64Var(
		(*uint64)(&args.thisServerId),
		"id",
		0,
		"server id of this server - this id should be in the cluster definition",
	)
	flag.Parse()

	// Check required args
	if args.cluster == "" {
		return nil, errors.New("flag -cluster is required")
	}
	if args.thisServerId == 0 {
		return nil, errors.New("flag -id is required and must not be 0")
	}

	return &args, nil
}
