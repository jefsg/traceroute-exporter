package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	traceroute "github.com/jefsg/traceroute-exporter/traceroute"
	log "github.com/sirupsen/logrus"
)

// test traceMiddleware
func TestEmptyTargetParam(t *testing.T) {
	resp := traceHandlerTestHelper("GET", "/trace?target=")
	assertEqual(t, resp.StatusCode, 400)

}

func TestNoTargetParam(t *testing.T) {
	resp := traceHandlerTestHelper("GET", "/trace")
	assertEqual(t, resp.StatusCode, 400)
}

func TestWithValidTarget(t *testing.T) {
	resp := traceHandlerTestHelper("GET", "/trace?target=google.com")
	assertEqual(t, resp.StatusCode, 200)
}

func TestWithFailedTrace(t *testing.T) {
	resp := traceHandlerTestHelper("GET", "/trace?target=fail.com")
	assertEqual(t, resp.StatusCode, 500)
}

// Helper functions

func traceHandlerTestHelper(method string, path string) *http.Response {
	req := httptest.NewRequest(method, "http://example.com"+path, nil)
	w := httptest.NewRecorder()
	traceHandler(w, req, fakeHTTPHandler(), fakeTracer)
	return w.Result()
}

func fakeHTTPHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success")) //nolint
	})
}

// fake tracer sends success response, unless host string contains the
//  substring "fail"
func fakeTracer(host string) ([]traceroute.Hop, error) {
	if strings.Contains(host, "fail") {
		return []traceroute.Hop{}, errors.New("contrived tracer error for testing")
	}
	return []traceroute.Hop{
		traceroute.Hop{
			Number:  "1",
			Name:    "hostname",
			Address: "address",
			Latency: 10.01,
		},
	}, nil
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func init() {
	// disable log output while testing
	log.SetOutput(ioutil.Discard)
}
