package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/utils/paginations"
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

func (s *instrumentingService) Create(ctx context.Context, req *model.User) (*model.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create").Add(1)
		s.requestLatency.With("method", "create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.Create(ctx, req)
}


func (s *instrumentingService) Update(ctx context.Context, user *model.User) (*model.User, error) {
	return s.next.Update(ctx, user)
}

func (s *instrumentingService) FindAll(pageable model.Pageable, result interface{}, where string, args ...interface{}) (*paginations.Pagination, error) {
	return s.next.FindAll(pageable, result, where, args...)
}
