package main

import (
	"net/http"
	"testing"

	"github.com/mf751/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	code, _, body := testServer.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
