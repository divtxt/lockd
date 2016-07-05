package locking

import (
	"encoding/json"
	"github.com/divtxt/raft"
)

//
type Cmd struct {
	Lock bool   `json:"a"`
	Name string `json:"n"`
}

//
func CmdSerialize(cmd *Cmd) (raft.Command, error) {
	jsonBytes, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

//
func CmdDeserialize(command raft.Command) (*Cmd, error) {
	var cmd *Cmd = &Cmd{}
	err := json.Unmarshal(command, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
