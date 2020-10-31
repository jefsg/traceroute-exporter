package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	port = "9805"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func hcHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("healthy")); err != nil {
		log.Error("error in path '/': " + err.Error())
	}
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	target, ok := r.URL.Query()["debug"]

	if !ok || len(target[0]) < 1 || target[0] != "on" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("healthy")); err != nil {
		log.Error("error in path '/debug': " + err.Error())
	}
}

func main() {

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server.
	http.Handle("/trace", traceMiddleware(promhttp.Handler()))
	http.HandleFunc("/", hcHandler)
	http.HandleFunc("/debug", debugHandler)

	log.Info("Listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
