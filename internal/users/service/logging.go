package service

import (
	"context"
	base "github.com/efrengarcial/framework/internal/platform/service"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type loggingService struct {
	logger log.Logger
	next   UserService
}


// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s UserService) UserService {
	return &loggingService{logger, s}
}

func (s *loggingService) Create(ctx context.Context, req *User) (u *User, err error) {
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

func (s *loggingService) Update(ctx context.Context, user *User) (*User, error) {
	return s.Update(ctx, user)
}


func (s *loggingService) FindAll(ctx context.Context, pageable *base.Pageable, result interface{}, where string, args ...interface{})(*base.Pagination, error) {
	return s.FindAll(ctx, pageable, result, where, args...)
}
