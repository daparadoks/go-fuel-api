//go:build e2e
// +build e2e

package tests

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndPoint(t *testing.T) {
	fmt.Println("running e2e test for health check endpoint")

	client := resty.New()
	resp, err := client.R().Get("http://localhost:3434/api/health")
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 200, resp.StatusCode())
}
