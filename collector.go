package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
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

func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		target, ok := r.URL.Query()["target"]

		if !ok || len(target[0]) < 1 {
			log.Error("target not provided in url params")
			// update success metric
			success.With(prometheus.Labels{"target": ""}).Set(-1)
		} else {
			if result, err := trace(target[0]); err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
					"host":  target[0],
				}).Error("traceroute error")
				success.With(prometheus.Labels{"target": target[0]}).Set(0)
			} else {
				log.WithFields(log.Fields{
					"host": target[0],
				}).Debug("traceroute succeeded")
				success.With(prometheus.Labels{"target": target[0]}).Set(1)
				hops.With(prometheus.Labels{"target": target[0]}).Set(float64(len(result)))
				for _, hop := range result {
					latency.With(prometheus.Labels{
						"target":      target[0],
						"hop_number":  string(hop.number),
						"hop_name":    hop.name,
						"hop_address": hop.address,
					}).Set(hop.latency)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
