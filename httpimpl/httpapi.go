package httpimpl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddLockApiEndpoints(e *gin.Engine, lockApi LockApi) {
	api := e.Group("/lock")
	api.GET("/:name", makeGetHandler(lockApi))
	api.POST("/:name", makeLockHandler(lockApi))
	api.DELETE("/:name", makeUnlockHandler(lockApi))
}

func makeGetHandler(lockApi LockApi) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		locked, _ := lockApi.IsLocked(name)
		if locked {
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		}
	}
}

func makeLockHandler(lockApi LockApi) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		commitChan := lockApi.Lock(name)
		if commitChan != nil {
			<-commitChan // FIXME: add timeout!
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Conflict"})
		}
	}
}

func makeUnlockHandler(lockApi LockApi) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		commitChan := lockApi.Unlock(name)
		if commitChan != nil {
			<-commitChan // FIXME: add timeout!
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		}
	}
}
