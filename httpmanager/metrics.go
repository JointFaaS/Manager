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
			Name:      "fn_processed_total",
			Help:      "Number of processed requests.",
		},
		[]string{"funcName"},
	)
	aliyunRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "manager",
			Subsystem: "requests",
			Name:      "fn_processed_aliyun",
			Help:      "Number of processed requests on aliyun.",
		},
		[]string{"funcName"},
	)
	workerRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "manager",
			Subsystem: "requests",
			Name:      "fn_processed_worker",
			Help:      "Number of processed requests on local VM worker.",
		},
		[]string{"funcName"},
	)
	priceMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "manager",
			Subsystem: "requests",
			Name:      "fn_processed_price",
			Help:      "Total price of HCloud.",
		},
		[]string{"funcName"},
	)
)

func (m *Manager) setMetrics() {
	prometheus.MustRegister(fnRequests)
	prometheus.MustRegister(aliyunRequests)
	prometheus.MustRegister(workerRequests)
	prometheus.MustRegister(priceMetrics)
	m.server.Handle("/metrics", promhttp.Handler())
}
