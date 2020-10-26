package metrics

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "github.com/techxmind/logserver/interface-defs"
	"github.com/techxmind/logserver/logger"
)

var (
	_hostname string

	// common action counter and gauge
	// labels: server, action
	_actionCounter metrics.Counter

	// outstanding request count
	// labels: server, method
	_outstandingRequest metrics.Gauge

	// labels: server, topic, app_type, event
	_eventCounter metrics.Counter

	// labels: server, method, code
	_requestCounter metrics.Counter
	// labels: server, method, code
	_requestDuration metrics.Histogram
)

func init() {
	if hostname, err := os.Hostname(); err != nil {
		logger.Error("Get hostname", "err", err)
	} else {
		_hostname = hostname
	}

	{
		opts := stdprometheus.CounterOpts{
			Namespace: "logserver",
			Help:      "logserver action counter",
			Name:      "action_count",
		}
		labels := []string{"server", "action"}
		_actionCounter = prometheus.NewCounterFrom(opts, labels)
	}

	{
		opts := stdprometheus.GaugeOpts{
			Namespace: "logserver",
			Help:      "logserver outstanding request",
			Name:      "outstanding",
		}
		labels := []string{"server", "method"}
		_outstandingRequest = prometheus.NewGaugeFrom(opts, labels)
	}

	{
		opts := stdprometheus.CounterOpts{
			Namespace: "logserver",
			Help:      "logserver event counter",
			Name:      "event_count",
		}
		labels := []string{"server", "topic", "app_type", "event"}
		_eventCounter = prometheus.NewCounterFrom(opts, labels)
	}

	{
		opts := stdprometheus.CounterOpts{
			Namespace: "logserver",
			Help:      "logserver request counter",
			Name:      "request_count",
		}
		labels := []string{"server", "method", "code"}
		_requestCounter = prometheus.NewCounterFrom(opts, labels)
	}

	{
		opts := stdprometheus.HistogramOpts{
			Namespace: "logserver",
			Help:      "logerver request duration",
			Name:      "request_duration",
			Buckets: []float64{
				.005, .01, .05, .1, .3, .5, .7, 1, 3, 5,
			},
		}
		labels := []string{"server", "method", "code"}
		_requestDuration = prometheus.NewHistogramFrom(opts, labels)
	}
}

// RequestMetrics is LabeledMiddleware, collect request count and latency
func RequestMetrics(label string, in endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		var code int32
		begin := time.Now()

		_outstandingRequest.With("server", _hostname, "method", label).Add(1)

		resp, err := in(ctx, req)

		if err != nil {
			code = 500
		} else {
			if r, ok := resp.(pb.Response); ok {
				code = r.Code
			}
		}

		AddRequest(label, code, time.Since(begin))

		_outstandingRequest.With("server", _hostname, "method", label).Add(-1)

		return resp, err
	}
}

// HttpHandler return prometheus http handler
func HttpHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		stdprometheus.DefaultRegisterer,
		promhttp.HandlerFor(stdprometheus.DefaultGatherer, promhttp.HandlerOpts{}))
}

func AddRequest(method string, code int32, du time.Duration) {
	codeStr := strconv.Itoa(int(code))

	_requestCounter.With(
		"server", _hostname,
		"method", method,
		"code", codeStr,
	).Add(1)

	_requestDuration.With(
		"server", _hostname,
		"method", method,
		"code", codeStr,
	).Observe(du.Seconds())
}

func CountEvent(topic, appType, event string) {
	_eventCounter.With(
		"server", _hostname,
		"topic", topic,
		"app_type", appType,
		"event", event,
	).Add(1)
}

func CounterAdd(action string, delta float64) {
	_actionCounter.With(
		"server", _hostname,
		"action", action,
	).Add(delta)
}
