package lockd_client

import (
	"fmt"
	"net/http"
)

type LockdClient struct {
	host string
	port uint
}

func NewLockdClient() *LockdClient {
	return &LockdClient{
		"127.0.0.1",
		2080,
	}
}

func (lc *LockdClient) IsLocked(name string) (bool, error) {
	return lc.lockish("GET", name, 404)
}

func (lc *LockdClient) Lock(name string) (bool, error) {
	return lc.lockish("POST", name, 409)
}

func (lc *LockdClient) Unlock(name string) (bool, error) {
	return lc.lockish("DELETE", name, 404)
}

func (lc *LockdClient) lockish(method string, name string, falseCode int) (bool, error) {
	// FIXME: does name need to escaped here?
	url := fmt.Sprintf("http://%s:%d/lock/%s", lc.host, lc.port, name)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	err = resp.Body.Close()
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	}
	if resp.StatusCode == falseCode {
		return false, nil
	}
	return false, fmt.Errorf("Unexpected response: %s", resp.Status)
}
