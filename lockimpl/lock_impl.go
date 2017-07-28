package lockimpl

import (
	"log"
	"time"

	"github.com/divtxt/lockd/backend"
	"github.com/divtxt/lockd/raftlock"
	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft"
	raft_config "github.com/divtxt/raft/config"
	raft_impl "github.com/divtxt/raft/impl"
	raft_log "github.com/divtxt/raft/log"
	raft_rps "github.com/divtxt/raft/rps"
)

const (
	MaxEntriesPerAppendEntry = 16

	TickerDuration     = 30 * time.Millisecond
	ElectionTimeoutLow = 150 * time.Millisecond
)

// Implementation of locking.LockApi
func NewLockApiImpl(
	cd util.ClusterDefinition,
	thisServerId raft.ServerId,
	logger *log.Logger,
) (raftlock.RaftLock, raft.IConsensusModule, error) {

	// --  Prepare raft ConsensusModule parameters

	raftPersistentState := raft_rps.NewIMPSWithCurrentTerm(0)

	raftLog := raft_log.NewInMemoryLog()

	timeSettings := raft_config.TimeSettings{TickerDuration, ElectionTimeoutLow}

	clusterInfo, err := raft_config.NewClusterInfo(cd.GetAllServerIds(), thisServerId)
	if err != nil {
		return nil, nil, err
	}

	rpcService := NewJsonRaftRpcService(cd, thisServerId, logger)

	// -- Make the LockBackend

	var lockBackend backend.LockBackend = backend.NewInMemoryBackend(0)

	// -- Make the RaftLock and raft ConsensusModule

	stateMachine, partialRaftLock := raftlock.NewRaftLock(lockBackend)

	raftCm, err := raft_impl.NewConsensusModule(
		raftPersistentState,
		raftLog,
		stateMachine,
		rpcService,
		clusterInfo,
		MaxEntriesPerAppendEntry,
		timeSettings,
		logger,
	)
	if err != nil {
		return nil, nil, err
	}

	raftLock := partialRaftLock.Finish(raftCm)

	return raftLock, raftCm, nil
}
