package util_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft"
)

func TestClusterDefinition(t *testing.T) {
	cd, err := util.LoadClusterDefinition("./cluster_test.json")
	if err != nil {
		t.Fatal(err)
	}

	asids := cd.GetAllServerIds()
	sort.Slice(asids, func(i, j int) bool { return asids[i] < asids[j] })

	if !reflect.DeepEqual(asids, []raft.ServerId{101, 102, 103}) {
		t.Error(asids)
	}

	if cd.GetHostPort(102) != "10.0.0.22:2081" {
		t.Error()
	}
}
