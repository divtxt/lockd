package lock

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AddEndpoints(e *gin.Engine) {
	e.POST("/lock", lockHandler)
	e.POST("/unlock", unlockHandler)
}

type LockRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func lockHandler(c *gin.Context) {
	var lockRequest LockRequest
	if err := c.Bind(&lockRequest); err == nil {
		log.Printf("Attempt to lock: %q", lockRequest.Name)
		c.JSON(http.StatusNotImplemented, gin.H{"status": "unimplemented"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
	}
}

func unlockHandler(c *gin.Context) {
	var lockRequest LockRequest
	if err := c.Bind(&lockRequest); err == nil {
		log.Printf("Attempt to unlock: %q", lockRequest.Name)
		c.JSON(http.StatusNotImplemented, gin.H{"status": "unimplemented"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
	}
}
