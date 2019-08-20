package metrics

import (
	"github.com/rcrowley/go-metrics"
)

const(
	HISTOGRAM_RESERVOIR_SIZE = 100
	HISTOGRAM_RESERVOIR_ALPHA = 0.015
)

func DeregisterAll() {
	metrics.Each(func(s string, i interface{}) {
		metrics.Unregister(s)
	})
}

// NewCounter create a new metrics Counter
func NewCounter(name string) metrics.Counter {
	return metrics.GetOrRegisterCounter(name, metrics.DefaultRegistry)
}

// NewMeter create a new metrics Meter
func NewMeter(name string) metrics.Meter {
	return metrics.GetOrRegisterMeter(name, metrics.DefaultRegistry)
}

// NewTimer create a new metrics Timer
func NewTimer(name string) metrics.Timer {
	return metrics.GetOrRegisterTimer(name, metrics.DefaultRegistry)
}

// NewGauge create a new metrics Gauge
func NewGauge(name string) metrics.Gauge {
	return metrics.GetOrRegisterGauge(name, metrics.DefaultRegistry)
}

func NewHistogram(name string) metrics.Histogram{
	return metrics.NewRegisteredHistogram(name,metrics.DefaultRegistry,metrics.NewExpDecaySample(HISTOGRAM_RESERVOIR_SIZE,HISTOGRAM_RESERVOIR_ALPHA))
}
