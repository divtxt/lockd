package locking

import (
	"time"

	"github.com/divtxt/lockd/raftlock"
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
func NewLockApiImpl() (*raftlock.RaftLock, error) {

	// --  Prepare raft ConsensusModule parameters

	raftPersistentState := raft_rps.NewIMPSWithCurrentTerm(0)

	raftLog := raft_log.NewInMemoryLog()

	timeSettings := raft_config.TimeSettings{TickerDuration, ElectionTimeoutLow}

	clusterInfo, err := raft_config.NewClusterInfo([]raft.ServerId{"_SOLO_"}, "_SOLO_")
	if err != nil {
		return nil, err
	}

	// -- Make the LockApi

	raftLock := raftlock.NewRaftLock(
		raftLog,
		[]string{}, // no initial locks
		0,          // initialCommitIndex
	)

	// -- Create the raft ConsensusModule
	raftCm, err := raft_impl.NewConsensusModule(
		raftPersistentState,
		raftLog,
		raftLock,
		nil, // should not actually need RpcService for single-node
		clusterInfo,
		MaxEntriesPerAppendEntry,
		timeSettings,
	)
	if err != nil {
		return nil, err
	}

	// Give RaftLock the IConsensusModule_AppendCommandOnly reference.
	raftLock.SetICMACO(raftCm)

	return raftLock, nil
}
