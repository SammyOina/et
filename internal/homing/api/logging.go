package api

import (
	"context"
	"fmt"
	"time"

	"github.com/mainflux/et/internal/homing"
	"github.com/mainflux/mainflux/logger"
)

var _ homing.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	hommingLogger logger.Logger
	svc           homing.Service
}

// LoggingMiddleware is a middleware that adds logging facilities to the core homing service.
func LoggingMiddleware(svc homing.Service, logger logger.Logger) homing.Service {
	return &loggingMiddleware{logger, svc}
}

// GetAll adds logging middleware to get all service.
func (lm *loggingMiddleware) GetAll(ctx context.Context, repo, token string, pm homing.PageMetadata) (telemetryPage homing.TelemetryPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method get all telemetry with token %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.hommingLogger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.hommingLogger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.GetAll(ctx, repo, token, pm)
}

// Save adds logging middleware to save service.
func (lm *loggingMiddleware) Save(ctx context.Context, t homing.Telemetry) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method save telemetry event took %s to complete", time.Since(begin))
		if err != nil {
			lm.hommingLogger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.hommingLogger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Save(ctx, t)
}
