package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mf751/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := &application{
		errLog:  log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
	}
	testServer := httptest.NewTLSServer(app.routes())
	defer testServer.Close()

	request, err := testServer.Client().Get(testServer.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, request.StatusCode, http.StatusOK)
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
