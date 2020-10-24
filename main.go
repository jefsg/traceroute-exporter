package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	port string = "8080"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/trace", traceMiddleware(promhttp.Handler()))
	log.Info("Listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
