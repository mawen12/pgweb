package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 指标统计
var (
	// sessions count
	sessionsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pgweb_sessions_count",
		Help: "Total number of database sessions",
	})

	// queries count
	queriesCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pgweb_queries_count",
		Help: "Total number of custom queries executed",
	})

	// healthy status
	healthyGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pgweb_healthy",
		Help: "Server health status",
	})

	// process start time
	startTimeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pgweb_process_start_time",
		Help: "Server start time, seconds since unix epoch",
	})

	// uptime
	uptimeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pgweb_uptime",
		Help: "Server application uptime in seconds",
	})
)

func init() {
	startTimeGauge.Set(float64(time.Now().Unix()))
}

func IncrementQueriesCount() {
	queriesCounter.Inc()
}

func SetSessionsCount(val int) {
	sessionsGauge.Set(float64(val))
}

func SetHealthy(val bool) {
	healthy := 0.0
	if val {
		healthy = 1.0
	}
	healthyGauge.Set(float64(healthy))
}
