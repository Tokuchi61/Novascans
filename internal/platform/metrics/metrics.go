package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Registry struct {
	mu                 sync.Mutex
	requestsTotal      uint64
	requestStatusTotal map[int]uint64
	durationCount      uint64
	durationSumSeconds float64
}

func NewRegistry() *Registry {
	return &Registry{
		requestStatusTotal: make(map[int]uint64),
	}
}

func (registry *Registry) ObserveRequest(statusCode int, duration time.Duration) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.requestsTotal++
	registry.requestStatusTotal[statusCode]++
	registry.durationCount++
	registry.durationSumSeconds += duration.Seconds()
}

func (registry *Registry) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	registry.mu.Lock()
	requestsTotal := registry.requestsTotal
	durationCount := registry.durationCount
	durationSumSeconds := registry.durationSumSeconds
	statusTotals := make(map[int]uint64, len(registry.requestStatusTotal))
	for statusCode, total := range registry.requestStatusTotal {
		statusTotals[statusCode] = total
	}
	registry.mu.Unlock()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	var builder strings.Builder
	builder.WriteString("# HELP novascans_http_requests_total Total number of HTTP requests.\n")
	builder.WriteString("# TYPE novascans_http_requests_total counter\n")
	builder.WriteString(fmt.Sprintf("novascans_http_requests_total %d\n", requestsTotal))
	builder.WriteString("# HELP novascans_http_requests_by_status_total Total number of HTTP requests grouped by status code.\n")
	builder.WriteString("# TYPE novascans_http_requests_by_status_total counter\n")

	statusCodes := make([]int, 0, len(statusTotals))
	for statusCode := range statusTotals {
		statusCodes = append(statusCodes, statusCode)
	}
	sort.Ints(statusCodes)

	for _, statusCode := range statusCodes {
		builder.WriteString(
			fmt.Sprintf(
				"novascans_http_requests_by_status_total{status=\"%d\"} %d\n",
				statusCode,
				statusTotals[statusCode],
			),
		)
	}

	builder.WriteString("# HELP novascans_http_request_duration_seconds_sum Total duration of HTTP requests in seconds.\n")
	builder.WriteString("# TYPE novascans_http_request_duration_seconds_sum counter\n")
	builder.WriteString(fmt.Sprintf("novascans_http_request_duration_seconds_sum %f\n", durationSumSeconds))
	builder.WriteString("# HELP novascans_http_request_duration_seconds_count Total number of observed HTTP request durations.\n")
	builder.WriteString("# TYPE novascans_http_request_duration_seconds_count counter\n")
	builder.WriteString(fmt.Sprintf("novascans_http_request_duration_seconds_count %d\n", durationCount))

	_, _ = w.Write([]byte(builder.String()))
}
