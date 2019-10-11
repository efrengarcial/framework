package user

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingService struct {
	logger log.Logger
	next   Service
}


// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Create(ctx context.Context, req *domain.User) (err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Create",
			"login", req.Login,
			"email", req.Email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Create(ctx, req)
}

func (s *loggingService) Update(ctx context.Context, user *domain.User) error {
	return s.Update(ctx, user)
}


func (s *loggingService) FindAll(ctx context.Context, pageable *domain.Pageable, result interface{}, where string, args ...interface{})(*domain.Pagination, error) {
	return s.FindAll(ctx, pageable, result, where, args...)
}
