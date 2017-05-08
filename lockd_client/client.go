package lockd_client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LockdClient struct {
	host string
	port uint
}

func NewLockdClient() *LockdClient {
	return &LockdClient{
		"127.0.0.1",
		2081,
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
	escapedName := url.PathEscape(name)
	lockUrl := fmt.Sprintf("http://%s:%d/lock/%s", lc.host, lc.port, escapedName)
	client := &http.Client{}
	req, err := http.NewRequest(method, lockUrl, nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	err2 := resp.Body.Close()
	if err != nil {
		return false, err
	}
	if err2 != nil {
		return false, err2
	}
	if resp.StatusCode == 200 {
		return true, nil
	}
	if resp.StatusCode == falseCode {
		return false, nil
	}
	if resp.StatusCode == 400 {
		return false, fmt.Errorf("Bad Request: %s", body)
	}
	return false, fmt.Errorf("Unexpected response: %s for %s /lock/%s", resp.Status, method, escapedName)
}
