package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const templateMetrics string = `Active connections: %d
server accepts handled requests
%d %d %d
Reading: %d Writing: %d Waiting: %d
`

// StubStatus represents NGINX stub_status metrics.
type StubStatus struct {
	Connections StubConnections
	Requests    int64
}

// StubConnections represents connections related metrics.
type StubConnections struct {
	Active   int64
	Accepted int64
	Handled  int64
	Reading  int64
	Writing  int64
	Waiting  int64
}

func getStubStatus(client *http.Client, url string) (*StubStatus, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v response, got %v", http.StatusOK, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(body)
	status, err := parseStubStatus(r)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func parseStubStatus(r io.Reader) (*StubStatus, error) {
	var s StubStatus
	if _, err := fmt.Fscanf(
		r,
		templateMetrics,
		&s.Connections.Active,
		&s.Connections.Accepted,
		&s.Connections.Handled,
		&s.Requests,
		&s.Connections.Reading,
		&s.Connections.Writing,
		&s.Connections.Waiting,
	); err != nil {
		return nil, err
	}
	return &s, nil
}
