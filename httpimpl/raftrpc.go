package httpimpl

import (
	"net/http"
	"strconv"

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
		// Get from ServerId
		fromString, ok := c.GetQuery("from")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'from' query parameter"})
			return
		}
		from, err := parseServerId(fromString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' query parameter"})
			return
		}

		// Parse json
		var rpc raft.RpcAppendEntries
		err = c.BindJSON(&rpc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Process rpc
		var reply *raft.RpcAppendEntriesReply
		reply, err = cm.ProcessRpcAppendEntries(from, &rpc)
		if err != nil {
			panic("ConsensusModule is shutdown")
		}

		// Reply
		c.JSON(http.StatusOK, reply)
	}
}

func makeRequestVoteHandler(cm raft.IConsensusModule) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get from ServerId
		fromString, ok := c.GetQuery("from")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'from' query parameter"})
			return
		}
		from, err := parseServerId(fromString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' query parameter"})
			return
		}

		// Parse json
		var rpc raft.RpcRequestVote
		err = c.BindJSON(&rpc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Process rpc
		var reply *raft.RpcRequestVoteReply
		reply, err = cm.ProcessRpcRequestVote(from, &rpc)
		if err != nil {
			panic("ConsensusModule is shutdown")
		}

		// Reply
		c.JSON(http.StatusOK, reply)
	}
}

func parseServerId(s string) (raft.ServerId, error) {
	serverId, err := strconv.ParseUint(s, 10, 64)
	return raft.ServerId(serverId), err
}
