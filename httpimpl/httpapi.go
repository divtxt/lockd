package httpimpl

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddLockApiEndpoints(e *gin.Engine, lockApi LockApi) {
	api := e.Group("/api")
	api.GET("/lock", makeGetHandler(lockApi))
	api.POST("/lock", makeLockHandler(lockApi))
	api.POST("/unlock", makeUnlockHandler(lockApi))
}

type LockRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func makeGetHandler(lockApi LockApi) func(c *gin.Context) {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		names := q["name"]
		if len(names) == 1 {
			name := names[0]
			locked, _ := lockApi.IsLocked(name)
			if locked {
				c.JSON(http.StatusOK, gin.H{})
			} else {
				c.JSON(http.StatusNotFound, gin.H{})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "GET /api/lock?name=... needs exactly 1 name parameter"})
		}
	}
}

func makeLockHandler(lockApi LockApi) func(c *gin.Context) {
	return func(c *gin.Context) {
		var lockRequest LockRequest
		if err := c.Bind(&lockRequest); err == nil {
			commitChan := lockApi.Lock(lockRequest.Name)
			if commitChan != nil {
				<-commitChan // FIXME: add timeout!
				c.JSON(http.StatusOK, gin.H{})
			} else {
				c.JSON(http.StatusConflict, gin.H{})
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
			commitChan := lockApi.Unlock(lockRequest.Name)
			if commitChan != nil {
				<-commitChan // FIXME: add timeout!
				c.JSON(http.StatusOK, gin.H{})
			} else {
				c.JSON(http.StatusConflict, gin.H{})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
		}
	}
}
