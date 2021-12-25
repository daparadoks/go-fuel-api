//go:build e2e
// +build e2e

package tests

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetCunsumptions(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/api/consumptions")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}
