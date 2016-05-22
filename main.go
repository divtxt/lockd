package main

import (
	"flag"
	"github.com/divtxt/lockd/ginx"
	"github.com/divtxt/lockd/lockd"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Parse command line args
	listenAddrPtr := flag.String("listen", ":2080", "listen address")
	flag.Parse()

	// Reset standard log flags (undo Gin's settings)
	log.SetFlags(log.LstdFlags)

	// Disable Gin debug logging
	gin.SetMode(gin.ReleaseMode)

	// Instantiate a lock service
	var l lockd.LockApi = lockd.NewLockApiImpl(lockd.NewInMemoryLSP())

	// Configure http service
	r := gin.New()
	r.Use(ginx.StdLogLogger())
	r.Use(ginx.StdLogRecovery())

	lockd.AddLockApiEndpoints(r, l)

	// Run forever / till stopped
	log.Println("Starting server on address:", *listenAddrPtr)
	r.Run(*listenAddrPtr)
}
