package main

import (
	"flag"
	"log"

	"github.com/divtxt/lockd/ginx"
	"github.com/divtxt/lockd/httpimpl"
	"github.com/divtxt/lockd/lockimpl"
	"github.com/gin-gonic/gin"
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
	var l httpimpl.LockApi
	var err error
	l, err = lockimpl.NewLockApiImpl()
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
