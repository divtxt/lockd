package httpimpl

import (
	"net/http"

	"github.com/divtxt/lockd/raftlock"
	"github.com/divtxt/lockd/util"
	"github.com/gin-gonic/gin"
)

func AddLockApiEndpoints(e *gin.Engine, lockApi raftlock.RaftLock) {
	api := e.Group("/lock")
	api.GET("/:name", makeGetHandler(lockApi))
	api.POST("/:name", makeLockHandler(lockApi))
	api.DELETE("/:name", makeUnlockHandler(lockApi))
}

func makeGetHandler(lockApi raftlock.RaftLock) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		if e := util.IsValidLockName(name); e != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": e})
		} else {
			locked := lockApi.IsLocked(name)
			if locked {
				c.JSON(http.StatusOK, gin.H{})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			}
		}
	}
}

func makeLockHandler(lockApi raftlock.RaftLock) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		if e := util.IsValidLockName(name); e != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": e})
		} else {
			success, err := lockApi.Lock(name)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			} else {
				if success {
					c.JSON(http.StatusOK, gin.H{})
				} else {
					c.JSON(http.StatusConflict, gin.H{"error": "Lock conflict"})
				}
			}
		}
	}
}

func makeUnlockHandler(lockApi raftlock.RaftLock) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		if e := util.IsValidLockName(name); e != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": e})
		} else {
			success, err := lockApi.Unlock(name)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			} else {
				if success {
					c.JSON(http.StatusOK, gin.H{})
				} else {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
				}
			}
		}
	}
}
