package mock

import "net/http"

// Client is the mock client
type Client struct {
	GetFunc func(url string) (*http.Response, error)
}

//Get is the mock client's `Do` func
func (m *Client) Get(url string) (*http.Response, error) {
	return GetFunc(url)
}

var (
	GetFunc func(url string) (*http.Response, error)
)
