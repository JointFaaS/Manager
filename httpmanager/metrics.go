package httpmanager

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "manager",
			Subsystem: "requests",
			Name: "processed_total",
			Help: "Number of processed requests.",
		},
		[]string{"funcName"},
	)
)

func (m *Manager) setMetrics() {
	m.server.Handle("/metrics", promhttp.Handler())
}
