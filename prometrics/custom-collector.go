package prometrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type customCollector struct {
	counterDesc *prometheus.Desc
	_typo       prometheus.ValueType
	_getter     func(Labels) float64
	_labels     Labels
}

func newCustomCollector(
	name string,
	description string,
	labels Labels,
	getter func(Labels) float64,
) *customCollector {
	return &customCollector{
		counterDesc: prometheus.NewDesc(name, description, nil, prometheus.Labels(labels)),
		_typo:       prometheus.GaugeValue,
		_getter:     getter,
		_labels:     labels,
	}
}

func (c *customCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.counterDesc
}

func (c *customCollector) Collect(ch chan<- prometheus.Metric) {
	value := c._getter(c._labels)
	ch <- prometheus.MustNewConstMetric(
		c.counterDesc,
		c._typo,
		value,
	)
}
