package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/callhome/internal/homing"
)

var _ homing.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     homing.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc homing.Service, counter metrics.Counter, latency metrics.Histogram) homing.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

// GetAll add metrics middleware to get all service.
func (mm *metricsMiddleware) GetAll(ctx context.Context, repo, token string, pm homing.PageMetadata) (homing.TelemetryPage, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "get all").Add(1)
		mm.latency.With("method", "get all").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.GetAll(ctx, repo, token, pm)
}

// Save adds metrics middleware to save service.
func (mm *metricsMiddleware) Save(ctx context.Context, t homing.Telemetry) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "save").Add(1)
		mm.latency.With("method", "save").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Save(ctx, t)
}
