package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	temps = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "nginx_request_metrics",
			Help: "The real-time metrics of the nginx.",
			// Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"server_name", "endpoint", "status_code"},
	)

	requests_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nginx_requests_total",
			Help: "Number of requests",
		},
		[]string{"host", "endpoint", "status_code"},
	)

	request_latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "nginx_request_latency_seconds",
			Help: "Request latency",
		},
		[]string{"host", "endpoint", "status_code"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(temps)
}

func main() {
	// www.example.com,127.0.0.1:80,162,6242,1,1,1,0,0,0,0,10,1,10,1
	// temps.With(prometheus.Labels{"server_name": "localhost"}).Inc()
	temps.WithLabelValues("www.example.com", "/admin/", "200").Observe(22)
	// temps.With(prometheus.Labels{"code": "404", "method": "GET"}).Observe(22)

	reg := prometheus.NewRegistry()
	reg.MustRegister(temps)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	// http.Handle("/metrics", promhttp.Handler())
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
