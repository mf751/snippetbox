package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		errLog:  log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, handler http.Handler) *testServer {
	testServerInstance := httptest.NewTLSServer(handler)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	testServerInstance.Client().Jar = jar
	testServerInstance.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{testServerInstance}
}

func (testServerInstance *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	request, err := testServerInstance.Client().Get(testServerInstance.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return request.StatusCode, request.Header, string(body)
}
