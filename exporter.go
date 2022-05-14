package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func probeHandler(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	scheme := r.URL.Query().Get("scheme")
	cluster := r.URL.Query().Get("cluster")
	service := r.URL.Query().Get("service")
	port := r.URL.Query().Get("port")
	path := r.URL.Query().Get("path")
	if region == "" {
		http.Error(w, "Must set `region`.", http.StatusBadRequest)
		return
	}
	if scheme == "" {
		scheme = "http"
	}
	if cluster == "" {
		http.Error(w, "Must set `cluster`.", http.StatusBadRequest)
		return
	}
	if service == "" {
		service = cluster
	}
	if port == "" {
		port = "443"
	}
	if path == "" {
		path = "/stub_status"
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(NginxCollectorConstructor(region, scheme, cluster, service, port, path))
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/probe", probeHandler)
	log.Fatal(http.ListenAndServe(":9113", nil))
}
