package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/utils/paginations"
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

func (s *loggingService) Create(ctx context.Context, req *model.User) (u *model.User, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Create",
			"login", req.Login,
			"email", req.Email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, req)
}

func (s *loggingService) Update(ctx context.Context, user *model.User) (*model.User, error) {
	return s.next.Update(ctx, user)
}


func (s *loggingService) FindAll(pageable *model.Pageable, result interface{}, where string, args ...interface{})(*paginations.Pagination, error) {
	return s.next.FindAll(pageable, result, where, args...)
}
