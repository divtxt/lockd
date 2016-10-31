package locking

import (
	"github.com/divtxt/lockd/raftlock"
	"github.com/divtxt/raft"
	raft_config "github.com/divtxt/raft/config"
	raft_impl "github.com/divtxt/raft/impl"
	raft_log "github.com/divtxt/raft/log"
	raft_rps "github.com/divtxt/raft/rps"
	"time"
)

const (
	MaxEntriesPerAppendEntry = 16

	TickerDuration     = 30 * time.Millisecond
	ElectionTimeoutLow = 150 * time.Millisecond
)

// Implementation of locking.LockApi
func NewLockApiImpl() (*raftlock.RaftLock, error) {

	// --  Prepare raft ConsensusModule parameters

	raftPersistentState := raft_rps.NewIMPSWithCurrentTerm(0)

	raftLog := raft_log.NewInMemoryLog()

	timeSettings := raft_config.TimeSettings{TickerDuration, ElectionTimeoutLow}

	clusterInfo, err := raft_config.NewClusterInfo([]raft.ServerId{"_SOLO_"}, "_SOLO_")
	if err != nil {
		return nil, err
	}

	// -- Create the raft ConsensusModule
	raftCm, err := raft_impl.NewConsensusModule(
		raftPersistentState,
		raftLog,
		nil, // should not actually need RpcService for single-node
		clusterInfo,
		MaxEntriesPerAppendEntry,
		timeSettings,
	)
	if err != nil {
		return nil, err
	}

	// -- Make the LockApi

	raftLock := raftlock.NewRaftLock(
		raftCm,
		raftLog,
		[]string{}, // no initial locks
		0,          // initialCommitIndex
	)

	raftCm.Start(raftLock)

	return raftLock, nil
}
