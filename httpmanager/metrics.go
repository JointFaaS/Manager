package httpmanager

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fnRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "manager",
			Subsystem: "requests",
			Name: "fn_processed_total",
			Help: "Number of processed requests.",
		},
		[]string{"funcName"},
	)
)

func (m *Manager) setMetrics() {
	prometheus.MustRegister(fnRequests)
	m.server.Handle("/metrics", promhttp.Handler())
}
