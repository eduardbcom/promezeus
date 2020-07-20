package prometrics

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Labels is an alias for prometheus.Labels value
type Labels prometheus.Labels

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
func Register(
	metricName string,
	labels []string,
	help string,
) {
	if _, ok := _metrics.Load(metricName); ok {
		return
	}

	// TODO: add support for counter as well
	m := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: help,
	}, labels)

	prometheus.MustRegister(m)

	_metrics.Store(metricName, m)
}

// RegisterGauge register and store collector for particular [metricKey, metricName] pair.
// collector function is used to provide metric value
func RegisterWithCollector(
	metricName string,
	labels Labels,
	help string,
	collector func(Labels) float64,
) {
	labelsMD5 := mapMD5(labels)
	k := generateCompoundKey(metricName, labelsMD5)

	if _, ok := _metrics.Load(k); ok {
		return
	}

	customCollector := newCustomCollector(
		metricName,
		help,
		labels,
		collector,
	)
	prometheus.MustRegister(customCollector)

	_metrics.Store(k, customCollector)
}

// Unregister unregisters metric fo specified name
func Unregister(metricName string) {
	metric, ok := _metrics.Load(metricName)
	if !ok {
		return
	}

	prometheus.Unregister(metric.(prometheus.Collector))

	_metrics.Delete(metricName)
}

// Unregister unregisters metric fo specified name
func UnregisterWithCollector(metricName string, labels Labels) {
	labelsMD5 := mapMD5(labels)
	k := generateCompoundKey(metricName, labelsMD5)

	metric, ok := _metrics.Load(k)
	if !ok {
		return
	}

	prometheus.Unregister(metric.(prometheus.Collector))

	_metrics.Delete(k)
}

// Increments metric with labels for one
func Inc(metricName string, labels Labels) {
	metric, ok := _metrics.Load(metricName)
	if !ok {
		panic("Inc on incorrect metric")
	}

	// TODO: add support for counter as well
	c := metric.(*prometheus.GaugeVec)

	c.With(prometheus.Labels(labels)).Inc()
}

func generateCompoundKey(metricKey, metricName string) string {
	return metricKey + metricName
}

func mapMD5(m map[string]string) string {
	if m == nil {
		return ""
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	b := new(bytes.Buffer)
	for _, key := range keys {
		fmt.Fprintf(b, "%s=%s", key, m[key])
	}

	h := md5.New()
	return string(h.Sum([]byte(b.String())))
}
