package prometrics

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Labels is an alias for prometheus.Labels value
type Labels prometheus.Labels

// GaugeType is an alias for prometheus.GaugeValue value
const GaugeType = prometheus.GaugeValue

var (
	_metrics = &sync.Map{}
)

// Server is a type for metrics server instances
type Server struct {
	srv *http.Server
}

// New returns new instance of the metrics server
func New() (*Server, error) {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:           "0.0.0.0:9100",
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server := &Server{
		srv: srv,
	}

	return server, nil
}

// Listen starts metrics server
func (r *Server) Listen() error {
	go func() {
		if err := r.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

// StopListen stops metrics server
func (r *Server) StopListen() error {
	return r.srv.Shutdown(context.Background())
}

// RegisterGauge register and store collector for particular [metricKey, metricName] pair.
// getter function is used to provide metric value
func RegisterGauge(
	metricKey string,
	typo prometheus.ValueType,
	metricName string,
	help string,
	labels Labels,
	getter func(Labels) float64,
) {
	k := generateCompoundKey(metricKey, metricName)

	if _, ok := _metrics.Load(k); ok {
		return
	}

	collector := newCustomCollector(
		typo,
		metricName,
		help,
		labels,
		getter,
	)
	prometheus.MustRegister(collector)

	_metrics.Store(k, collector)
}

// Unregister unregisters metric fo specified name
func Unregister(metricKey, metricName string) {
	k := generateCompoundKey(metricKey, metricName)

	metric, ok := _metrics.Load(k)
	if !ok {
		return
	}

	prometheus.Unregister(metric.(prometheus.Collector))

	_metrics.Delete(k)
}

func generateCompoundKey(metricKey, metricName string) string {
	return metricKey + metricName
}
