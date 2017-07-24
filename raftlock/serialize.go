package raftlock

import (
	"encoding/json"

	"github.com/divtxt/raft"
)

type lockAction struct {
	Lock bool   `json:"a"`
	Name string `json:"n"`
}

func MakeLockCommand(name string) (raft.Command, error) {
	return json.Marshal(&lockAction{true, name})
}

func MakeUnlockCommand(name string) (raft.Command, error) {
	return json.Marshal(&lockAction{false, name})
}

func lockActionDeserialize(command raft.Command) (*lockAction, error) {
	var cmd *lockAction = &lockAction{}
	err := json.Unmarshal(command, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
