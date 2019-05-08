package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/mocks"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

var logger log.Logger


func init() {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger,
		"svc", "users",
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
}

//https://github.com/bxcodec/go-clean-arch/blob/master/article/usecase/article_ucase_test.go
func TestGet(t *testing.T) {
	mockRepository :=  new(mocks.Repository)

	mockUser := &model.User {
		FirstName: "Efren",
		LastName: "Garcia",
		Email: "efren.gl@gmail.com",
		Login: "efren",
	}

	t.Run("success", func(t *testing.T) {
		mockRepository.On("Insert", mock.AnythingOfType("*model.User")).Return(mockUser, nil).Once()
		u := NewService(mockRepository,  logger)

		user, err := u.Create(context.TODO(), mockUser)

		assert.NoError(t, err)
		assert.NotNil(t, user)

		mockRepository.AssertExpectations(t)
	})
}