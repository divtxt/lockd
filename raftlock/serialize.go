package raftlock

import (
	"encoding/json"

	"github.com/divtxt/raft"
)

type lockAction struct {
	Lock bool   `json:"a"`
	Name string `json:"n"`
}

type lockActionResult struct {
	success bool
}

func lockActionSerialize(cmd *lockAction) (raft.Command, error) {
	jsonBytes, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func lockActionDeserialize(command raft.Command) (*lockAction, error) {
	var cmd *lockAction = &lockAction{}
	err := json.Unmarshal(command, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
