package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	ecsServiceReachable   = 1
	ecsServiceUnreachable = 0
	nginxUp               = 1
	nginxDown             = 0
)

type NginxCollector struct {
	region  string
	scheme  string
	cluster string
	service string
	port    string
	path    string
	metrics map[string]*prometheus.Desc
}

func NginxCollectorConstructor(region string, scheme string, cluster string, service string, port string, path string) *NginxCollector {
	ecsLevelLabels := []string{"cluster", "service"}
	nginxLevelLabels := []string{"cluster", "service", "target_ip"}
	return &NginxCollector{
		region:  region,
		scheme:  scheme,
		cluster: cluster,
		service: service,
		port:    port,
		path:    path,
		metrics: map[string]*prometheus.Desc{
			"ecsServiceReachable": prometheus.NewDesc("nginx_ecs_service_reachable", "Find the ecs service", ecsLevelLabels, nil),
			"nginxUp":             prometheus.NewDesc("nginx_up", "Status of the last metric scrape", nginxLevelLabels, nil),
			"connectionsActive":   prometheus.NewDesc("nginx_connections_active", "Active client connections", nginxLevelLabels, nil),
			"connectionsAccepted": prometheus.NewDesc("nginx_connections_accepted", "Accepted client connections", nginxLevelLabels, nil),
			"connectionsHandled":  prometheus.NewDesc("nginx_connections_handled", "Handled client connections", nginxLevelLabels, nil),
			"connectionsReading":  prometheus.NewDesc("nginx_connections_reading", "Connections where NGINX is reading the request header", nginxLevelLabels, nil),
			"connectionsWriting":  prometheus.NewDesc("nginx_connections_writing", "Connections where NGINX is writing the response back to the client", nginxLevelLabels, nil),
			"connectionsWaiting":  prometheus.NewDesc("nginx_connections_waiting", "Idle client connections", nginxLevelLabels, nil),
			"httpRequestsTotal":   prometheus.NewDesc("nginx_http_requests_total", "Total http requests", nginxLevelLabels, nil),
		},
	}
}

func (nc *NginxCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range nc.metrics {
		ch <- m
	}
}

func (nc *NginxCollector) Collect(ch chan<- prometheus.Metric) {
	targetIPs, err := getTargetIPs(nc.region, nc.cluster, nc.service)
	if err != nil {
		log.Print(err)
		ch <- prometheus.MustNewConstMetric(nc.metrics["ecsServiceReachable"], prometheus.GaugeValue, float64(ecsServiceUnreachable), nc.cluster, nc.service)
		return
	}
	ch <- prometheus.MustNewConstMetric(nc.metrics["ecsServiceReachable"], prometheus.GaugeValue, float64(ecsServiceReachable), nc.cluster, nc.service)

	probeClient := http.DefaultClient
	for _, ip := range targetIPs {
		url := fmt.Sprintf("%s://%s:%s%s", nc.scheme, ip, nc.port, nc.path)
		status, err := getStubStatus(probeClient, url)
		if err != nil {
			log.Print(err)
			ch <- prometheus.MustNewConstMetric(nc.metrics["nginxUp"], prometheus.GaugeValue, float64(nginxDown), nc.cluster, nc.service, ip)
		} else {
			ch <- prometheus.MustNewConstMetric(nc.metrics["nginxUp"], prometheus.GaugeValue, float64(nginxUp), nc.cluster, nc.service, ip)

			ch <- prometheus.MustNewConstMetric(nc.metrics["connectionsActive"], prometheus.GaugeValue, float64(status.Connections.Active), nc.cluster, nc.service, ip)
			ch <- prometheus.MustNewConstMetric(nc.metrics["connectionsAccepted"], prometheus.CounterValue, float64(status.Connections.Accepted), nc.cluster, nc.service, ip)
			ch <- prometheus.MustNewConstMetric(nc.metrics["connectionsHandled"], prometheus.CounterValue, float64(status.Connections.Handled), nc.cluster, nc.service, ip)
			ch <- prometheus.MustNewConstMetric(nc.metrics["connectionsReading"], prometheus.GaugeValue, float64(status.Connections.Reading), nc.cluster, nc.service, ip)
			ch <- prometheus.MustNewConstMetric(nc.metrics["connectionsWriting"], prometheus.GaugeValue, float64(status.Connections.Writing), nc.cluster, nc.service, ip)
			ch <- prometheus.MustNewConstMetric(nc.metrics["connectionsWaiting"], prometheus.GaugeValue, float64(status.Connections.Waiting), nc.cluster, nc.service, ip)
			ch <- prometheus.MustNewConstMetric(nc.metrics["httpRequestsTotal"], prometheus.CounterValue, float64(status.Requests), nc.cluster, nc.service, ip)
		}
	}
}
