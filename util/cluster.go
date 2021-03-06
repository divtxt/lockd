package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/divtxt/raft"
)

// ClusterDefinition holds the server ids and addresses of the servers in the cluster.
type ClusterDefinition map[raft.ServerId]string

// LoadClusterDefinition loads the cluster definition from given json file.
//
// The json should be of the form: {<server id>: "host:port", ...}
// Example: {1: "lockd1:2081", 2: "lockd2:2082", 3: "lockd3:2083"}
//
// The json should be of the form: {"server-id": "host:port", ...}
// server ids should be positive integers, but as strings since json keys must be strings
// example: {\"1\": \"lockd1:2081\", \"2\": \"lockd2:2082\", \"3\": \"lockd3:2083\"}
func LoadClusterDefinition(name string) (ClusterDefinition, error) {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	m := map[raft.ServerId]string{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (cd ClusterDefinition) GetAllServerIds() []raft.ServerId {
	var asids []raft.ServerId
	for k := range cd {
		asids = append(asids, k)
	}

	return asids
}

func (cd ClusterDefinition) GetHostPort(sid raft.ServerId) string {
	hostPort, ok := cd[sid]
	if !ok {
		// FIXME: change to error
		panic(fmt.Sprintf("unknown server id: %v", sid))
	}
	return hostPort
}
