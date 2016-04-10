package lockd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AddLockApiEndpoints(e *gin.Engine, lockApi LockApi) {
	e.POST("/lock", makeLockHandler(lockApi))
	e.POST("/unlock", makeUnlockHandler(lockApi))
}

type LockRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func makeLockHandler(lockApi LockApi) func(c *gin.Context) {
	return func(c *gin.Context) {
		var lockRequest LockRequest
		if err := c.Bind(&lockRequest); err == nil {
			log.Printf("Attempt to lock: %q", lockRequest.Name)
			success, err := lockApi.Lock(lockRequest.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{})
			} else {
				if success {
					c.JSON(http.StatusOK, gin.H{})
				} else {
					c.JSON(http.StatusConflict, gin.H{})
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
		}
	}
}

func makeUnlockHandler(lockApi LockApi) func(c *gin.Context) {
	return func(c *gin.Context) {
		var lockRequest LockRequest
		if err := c.Bind(&lockRequest); err == nil {
			log.Printf("Attempt to unlock: %q", lockRequest.Name)
			success, err := lockApi.Unlock(lockRequest.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{})
			} else {
				if success {
					c.JSON(http.StatusOK, gin.H{})
				} else {
					c.JSON(http.StatusConflict, gin.H{})
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
		}
	}
}
