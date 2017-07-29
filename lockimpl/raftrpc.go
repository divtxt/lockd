package lockimpl

import (
	"log"
	"sync"

	"fmt"

	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft"
)

type JsonRaftRpcService struct {
	cd           util.ClusterDefinition
	thisServerId raft.ServerId
	logger       *log.Logger

	// mutable fields - use mutex for access!
	mutex             sync.Mutex
	_toServerErrCount map[raft.ServerId]uint // consecutive err count
}

func NewJsonRaftRpcService(
	cd util.ClusterDefinition,
	thisServerId raft.ServerId,
	logger *log.Logger,
) *JsonRaftRpcService {
	tsec := make(map[raft.ServerId]uint)
	for _, sid := range cd.GetAllServerIds() {
		tsec[sid] = 0
	}
	return &JsonRaftRpcService{
		cd,
		thisServerId,
		logger,
		sync.Mutex{},
		tsec,
	}
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
		jrrs.logRpcError(toServer, err)
		return nil
	}
	jrrs.logRpcSuccess(toServer)
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
		jrrs.logRpcError(toServer, err)
		return nil
	}
	jrrs.logRpcSuccess(toServer)
	return &reply
}

func (jrrs *JsonRaftRpcService) logRpcSuccess(toServer raft.ServerId) {
	jrrs.mutex.Lock()
	defer jrrs.mutex.Unlock()
	if jrrs._toServerErrCount[toServer] > 0 {
		jrrs.logger.Println(
			"[lockd] rpc to", toServer,
			"- success after", jrrs._toServerErrCount[toServer], "consecutive errors",
		)
		jrrs._toServerErrCount[toServer] = 0
	}
}

func (jrrs *JsonRaftRpcService) logRpcError(toServer raft.ServerId, err error) {
	jrrs.mutex.Lock()
	defer jrrs.mutex.Unlock()
	if jrrs._toServerErrCount[toServer] == 0 {
		jrrs.logger.Println("[lockd] rpc to", toServer, "- error:", err)
		jrrs.logger.Println("[lockd] silencing further rpc errors to", toServer)
	}
	jrrs._toServerErrCount[toServer]++
}
