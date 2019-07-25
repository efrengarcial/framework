package service

import (
	"context"
	"github.com/efrengarcial/framework/internal/platform/database"
	base "github.com/efrengarcial/framework/internal/platform/service"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           UserService
}


// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s UserService) UserService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		next:           s,
	}
}

func (s *instrumentingService) Create(ctx context.Context, req *User) (*User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create").Add(1)
		s.requestLatency.With("method", "create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Create(ctx, req)
}


func (s *instrumentingService) Update(ctx context.Context, user *User) (*User, error) {
	return s.Update(ctx, user)
}

func (s *instrumentingService) FindAll(pageable *base.Pageable, result interface{}, where string, args ...interface{}) (*database.Pagination, error) {
	return s.FindAll(pageable, result, where, args...)
}
