package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	DatabaseQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "database_queries_total",
		Help: "Total number of database queries",
	}, []string{"type", "table"})

	DatabaseQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "database_query_duration_seconds",
		Help:    "Duration of database queries",
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
	}, []string{"type", "table"})

	AuthRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_requests_total",
		Help: "Total number of authentication requests",
	}, []string{"type", "status"})

	WishOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "wish_operations_total",
		Help: "Total number of wish operations",
	}, []string{"type", "status"})
)

func RecordDatabaseQuery(queryType, table string, duration float64) {
	DatabaseQueries.WithLabelValues(queryType, table).Inc()
	DatabaseQueryDuration.WithLabelValues(queryType, table).Observe(duration)
}

func RecordAuthRequest(requestType, status string) {
	AuthRequests.WithLabelValues(requestType, status).Inc()
}

func RecordWishOperation(operationType, status string) {
	WishOperations.WithLabelValues(operationType, status).Inc()
}

func Init() {
	promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_info",
		Help: "Application information",
		ConstLabels: prometheus.Labels{
			"version": "1.0.0",
		},
	}).Set(1)
}
