package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "streampulse_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "streampulse_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	ActiveStreams = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "streampulse_active_streams",
			Help: "Number of currently active streams",
		},
	)

	ActiveListeners = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "streampulse_active_listeners",
			Help: "Number of active listeners per stream",
		},
		[]string{"stream_id"},
	)

	StreamDisconnections = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "streampulse_stream_disconnections_total",
			Help: "Total number of abrupt stream disconnections (business metric)",
		},
		[]string{"stream_id", "reason"},
	)

	StreamBytesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "streampulse_stream_bytes_total",
			Help: "Total bytes broadcast per stream (business metric)",
		},
		[]string{"stream_id"},
	)

	AuthLoginsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "streampulse_auth_logins_total",
			Help: "Total login attempts, labelled by outcome",
		},
		[]string{"result"}, // success | invalid_credentials | error
	)
)
