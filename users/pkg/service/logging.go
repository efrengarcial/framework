package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/go-kit/kit/log"
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
		s.logger.Log(
			"method", "Create",
			"login", req.Login,
			"email", req.Email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, req)
}
