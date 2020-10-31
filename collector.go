package main

import (
	"net/http"

	traceroute "github.com/jefsg/traceroute-exporter/traceroute"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	latency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "route_hop_latency",
			Help: "Latency to the indicated hop (in seconds)",
		},
		[]string{"target", "hop_number", "hop_name", "hop_address"},
	)

	success = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "route_success",
			Help: "Indicates whether the trace was successful (1 = success, 0 = failure, -1 = exporter error)",
		},
		[]string{"target"},
	)

	hops = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "route_hop_count",
			Help: "Number of hops taken along route",
		},
		[]string{"target"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(latency)
	prometheus.MustRegister(success)
	prometheus.MustRegister(hops)
}

func traceHandler(w http.ResponseWriter, r *http.Request, next http.Handler, tracerFunc traceroute.Tracer) {

	log.WithFields(log.Fields{
		"method": r.Method,
	}).Debug(r.URL.String())

	target, ok := r.URL.Query()["target"]

	if !ok || len(target[0]) < 1 {
		errorMsg := "target not provided in url params"
		log.Error(errorMsg)
		// update success metric
		success.With(prometheus.Labels{"target": ""}).Set(-1)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(errorMsg)); err != nil {
			log.Error(err.Error())
		}
	} else if result, err := tracerFunc(target[0]); err != nil {
		errorMsg := "traceroute error"
		log.WithFields(log.Fields{
			"error": err.Error(),
			"host":  target[0],
		}).Error(errorMsg)
		success.With(prometheus.Labels{"target": target[0]}).Set(0)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(errorMsg)); err != nil {
			log.Error(err.Error())
		}
	} else {
		log.WithFields(log.Fields{
			"host": target[0],
		}).Debug("traceroute succeeded")
		success.With(prometheus.Labels{"target": target[0]}).Set(1)
		hops.With(prometheus.Labels{"target": target[0]}).Set(float64(len(result)))
		for _, hop := range result {
			latency.With(prometheus.Labels{
				"target":      target[0],
				"hop_number":  string(hop.Number),
				"hop_name":    hop.Name,
				"hop_address": hop.Address,
			}).Set(hop.Latency)
		}
	}
	next.ServeHTTP(w, r)
}

func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceHandler(w, r, next, traceroute.Trace)
	})
}
