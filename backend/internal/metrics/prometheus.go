package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Certificate metrics
	CertificatesTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificates_total",
			Help: "Total number of certificates",
		},
		[]string{"status"},
	)

	CertificatesExpiring = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificates_expiring",
			Help: "Number of certificates expiring",
		},
		[]string{"days"},
	)

	CertificateIssuanceDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "certificate_issuance_duration_seconds",
			Help:    "Time taken to issue a certificate",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
	)

	CertificateIssuanceTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "certificate_issuance_total",
			Help: "Total number of certificate issuances",
		},
		[]string{"adapter", "status"},
	)

	// Agent metrics
	AgentsTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "agents_total",
			Help: "Total number of agents",
		},
		[]string{"status"},
	)

	AgentCheckinDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "agent_checkin_duration_seconds",
			Help:    "Time taken for agent check-in",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Database metrics
	DatabaseConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"query_type"},
	)

	// Adapter metrics
	AdapterRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "adapter_requests_total",
			Help: "Total adapter requests",
		},
		[]string{"adapter", "status"},
	)

	AdapterRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "adapter_request_duration_seconds",
			Help:    "Adapter request duration",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"adapter"},
	)

	// Audit metrics
	AuditLogsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "audit_logs_total",
			Help: "Total audit log entries",
		},
		[]string{"action", "entity_type"},
	)
)

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration) {
	status := "2xx"
	if statusCode >= 400 && statusCode < 500 {
		status = "4xx"
	} else if statusCode >= 500 {
		status = "5xx"
	}

	HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordCertificateIssuance records certificate issuance metrics
func RecordCertificateIssuance(adapter string, success bool, duration time.Duration) {
	status := "success"
	if !success {
		status = "failure"
	}
	CertificateIssuanceTotal.WithLabelValues(adapter, status).Inc()
	CertificateIssuanceDuration.Observe(duration.Seconds())
}
