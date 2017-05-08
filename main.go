package main

import (
	"flag"
	"log"
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
	// Parse command line args
	listenAddrPtr := flag.String("listen", ":2081", "listen address")
	clusterPtr := flag.String(
		"cluster",
		"",
		"cluster definition file - json file that describes the lockd cluster\n"+
			"    \t    the json should be of the form: {\"server-id\": \"host:port\", ...}\n"+
			"    \t    server ids should be positive integers, but as strings since json keys must be strings\n"+
			"    \t    example: {\"1\": \"lockd1:2081\", \"2\": \"lockd2:2082\", \"3\": \"lockd3:2083\"}",
	)
	thisServerIdPtr := flag.Uint64(
		"id",
		0,
		"server id of this server - this id should be in the cluster definition",
	)
	flag.Parse()

	// Check required args
	if *clusterPtr == "" {
		fmt.Fprintln(os.Stderr, "Error: flag -cluster is required")
		flag.Usage()
		os.Exit(2)
	}
	if *thisServerIdPtr == 0 {
		fmt.Fprintln(os.Stderr, "Error: flag -id is required")
		flag.Usage()
		os.Exit(2)
	}
	// Reset standard log flags (undo Gin's settings)
	log.SetFlags(log.LstdFlags)

	// Disable Gin debug logging
	gin.SetMode(gin.ReleaseMode)

	// Process args
	log.Println("lockd initializing")
	cd, err := util.LoadClusterDefinition(*clusterPtr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loadig cluster definition: %v", err)
		os.Exit(2)
	}
	log.Println("cluster:", cd)
	log.Println("thisServerId:", *thisServerIdPtr)

	// Instantiate a lock service
	var l httpimpl.LockApi
	l, err = lockimpl.NewLockApiImpl(cd.GetAllServerIds(), raft.ServerId(*thisServerIdPtr))
	if err != nil {
		panic(err)
	}

	// Configure http service
	r := gin.New()
	r.Use(ginx.StdLogLogger())
	r.Use(ginx.StdLogRepanic())

	httpimpl.AddLockApiEndpoints(r, l)

	// Run forever / till stopped
	log.Println("Starting server on address:", *listenAddrPtr)
	err = r.Run(*listenAddrPtr)
	if err != nil {
		panic(err)
	}
}
