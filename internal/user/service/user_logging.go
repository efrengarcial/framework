package service

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"github.com/efrengarcial/framework/internal/user"
	log "github.com/sirupsen/logrus"
	"time"
)

type loggingService struct {
	logger     *log.Logger
	next   user.Service
}


// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger *log.Logger, s user.Service) *loggingService {
	return &loggingService{logger, s}
}

func (s *loggingService) Create(ctx context.Context, req *domain.User) (err error) {
	defer func(begin time.Time) {
		s.logger.WithFields(log.Fields{
			"method" : "Create",
			"login": req.Login,
			"email" : req.Email,
			"took" : time.Since(begin),
			"err" : err,
		}).Info("UserService", "Create")

	}(time.Now())
	return s.next.Create(ctx, req)
}

func (s *loggingService) Update(ctx context.Context, user *domain.User) error {
	return s.next.Update(ctx, user)
}


func (s *loggingService) FindAll(ctx context.Context, pageable *domain.Pageable, result interface{}, where string, args ...interface{})(*domain.Pagination, error) {
	return s.next.FindAll(ctx, pageable, result, where, args...)
}
