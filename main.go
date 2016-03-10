package main

import (
	"flag"
	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line args
	listenAddrPtr := flag.String("listen", ":8080", "listen address")
	flag.Parse()

	// Run http service
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello, world!",
		})
	})
	r.Run(*listenAddrPtr)
}
