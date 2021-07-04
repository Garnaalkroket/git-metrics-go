package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	/// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	gitMetrics := newGitCollector(regexp.MustCompile("CLONE RUN \\'\\d.\\d\\d\\' REPO \\'.*\\'"))
	prometheus.MustRegister(gitMetrics)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
