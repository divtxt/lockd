package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Parse command line args
	listenAddrPtr := flag.String("listen", ":8080", "listen address")
	flag.Parse()

	// Initialize logging
	log.SetFlags(log.LstdFlags)

	// Prevent Gin debug logging
	gin.SetMode(gin.ReleaseMode)

	// Run http service
	log.Println("Starting server on address:", *listenAddrPtr)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello, world!",
		})
	})
	r.Run(*listenAddrPtr)
}
