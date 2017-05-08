package httpimpl

import (
	"net/http"

	"github.com/divtxt/raft"
	"github.com/gin-gonic/gin"
)

func AddRaftRpcEndpoints(e *gin.Engine, cm raft.IConsensusModule) {
	api := e.Group("/raft")
	api.POST("/AppendEntries", makeAppendEntriesHandler(cm))
	api.POST("/RequestVote", makeRequestVoteHandler(cm))
}

func makeAppendEntriesHandler(cm raft.IConsensusModule) func(c *gin.Context) {
	return func(c *gin.Context) {
		var rpc raft.RpcAppendEntries
		err := c.BindJSON(&rpc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			var reply *raft.RpcAppendEntriesReply
			reply = cm.ProcessRpcAppendEntries(111, &rpc)
			if reply == nil {
				panic("ConsensusModule is shutdown")
			}
			c.JSON(http.StatusOK, reply)
		}
	}
}

func makeRequestVoteHandler(cm raft.IConsensusModule) func(c *gin.Context) {
	return func(c *gin.Context) {
		var rpc raft.RpcRequestVote
		err := c.BindJSON(&rpc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			var reply *raft.RpcRequestVoteReply
			reply = cm.ProcessRpcRequestVote(111, &rpc)
			if reply == nil {
				panic("ConsensusModule is shutdown")
			}
			c.JSON(http.StatusOK, reply)
		}
	}
}
