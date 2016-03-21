package main

import (
	"flag"
	api_lock "github.com/divtxt/lockd/api"
	"github.com/divtxt/lockd/misc"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Parse command line args
	listenAddrPtr := flag.String("listen", ":2080", "listen address")
	flag.Parse()

	// Initialize logging
	log.SetFlags(log.LstdFlags)

	// Prevent Gin debug logging
	gin.SetMode(gin.ReleaseMode)

	// Run http service
	log.Println("Starting server on address:", *listenAddrPtr)

	r := gin.New()
	r.Use(misc.StdLogLogger())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello, world!",
		})
	})

	api_lock.AddEndpoints(r)

	r.Run(*listenAddrPtr)
}
