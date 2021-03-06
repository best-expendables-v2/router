package middleware

import (
	"net/http"

	lgoprometheus "github.com/best-expendables-v2/router/internal/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Prometheus() func(http.Handler) http.Handler {
	counters := lgoprometheus.DefaultCounters()

	return func(h http.Handler) http.Handler {
		h = promhttp.InstrumentHandlerInFlight(counters.InFlightGauge, h)
		h = lgoprometheus.InstrumentHandlerCounter(counters.ReqTotalCounter, h)
		h = lgoprometheus.InstrumentHandlerDuration(counters.ReqDurationHistogram, h)
		return h
	}
}
