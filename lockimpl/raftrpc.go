package lockimpl

import (
	"log"

	"fmt"

	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft"
)

type JsonRaftRpcService struct {
	cd           util.ClusterDefinition
	thisServerId raft.ServerId
	logger       *log.Logger
}

func NewJsonRaftRpcService(
	cd util.ClusterDefinition,
	thisServerId raft.ServerId,
	logger *log.Logger,
) *JsonRaftRpcService {
	return &JsonRaftRpcService{cd, thisServerId, logger}
}

func (jrrs *JsonRaftRpcService) RpcAppendEntries(
	toServer raft.ServerId,
	rpc *raft.RpcAppendEntries,
) *raft.RpcAppendEntriesReply {
	hostPort := jrrs.cd.GetHostPort(toServer)
	url := fmt.Sprintf("http://%s/raft/AppendEntries?from=%d", hostPort, uint64(jrrs.thisServerId))
	var reply raft.RpcAppendEntriesReply
	err := util.JsonPost(url, rpc, &reply)
	if err != nil {
		jrrs.logger.Printf("[lockd] rpc to %v -> error: %v\n", toServer, err)
		return nil
	}
	return &reply
}

func (jrrs *JsonRaftRpcService) RpcRequestVote(
	toServer raft.ServerId,
	rpc *raft.RpcRequestVote,
) *raft.RpcRequestVoteReply {
	hostPort := jrrs.cd.GetHostPort(toServer)
	url := fmt.Sprintf("http://%s/raft/RequestVote?from=%d", hostPort, uint64(jrrs.thisServerId))
	var reply raft.RpcRequestVoteReply
	err := util.JsonPost(url, rpc, &reply)
	if err != nil {
		jrrs.logger.Printf("[lockd] rpc to %v -> error: %v\n", toServer, err)
		return nil
	}
	return &reply
}
