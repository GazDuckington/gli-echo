package utils

import (
	"io"
	"net/http"
	"time"
)

type Requests struct {
	Client *http.Client
	URL    string
}

type RequestsConfig struct {
	URL     string
	Timeout time.Duration
}

type Response struct {
	Message map[string][]string `json:"message"`
	Status  string              `json:"status"`
}
type GenericResponse struct {
	Message []string `json:"message"`
	Status  string   `json:"status"`
}

// request initializer
func NewRequests(config RequestsConfig) *Requests {
	return &Requests{
		Client: &http.Client{
			Timeout: config.Timeout,
		},
		URL: config.URL,
	}
}

// get
func (r *Requests) Get() ([]byte, error) {
	resp, err := r.Client.Get(r.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
