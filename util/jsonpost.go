package util

import (
	"bytes"
	"encoding/json"
	"net/http"

	"fmt"
)

// JsonPost performs a HTTP POST with json send and reply payloads.
//
// The json encoding of the given data is POSTed the given url and
// the reply is json is decoded into the given reply structure.
func JsonPost(url string, data interface{}, reply interface{}) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("JsonPost %v: unexpected http status code: %v", url, resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return err
	}
	return nil
}
