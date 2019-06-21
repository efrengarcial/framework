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
		"svc", "UserService",
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
}

//https://github.com/bxcodec/go-clean-arch/blob/master/article/usecase/article_ucase_test.go
// mockery -name=Repository
func TestInsert(t *testing.T) {
	mockUserRepository :=  new(mocks.UserRepository)

	mockAuthority := model.Authority{Model: model.Model{ID: 1}}

	mockUser := &model.User {
		FirstName: "Juan",
		LastName: "Perez",
		Email: "juan.perez@gmail.com",
		Login: "jperez",
		Authorities: []model.Authority{ mockAuthority },
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepository.On("Insert", mock.AnythingOfType("*model.User")).Return(mockUser, nil).Once()
		u := NewService(mockUserRepository,  logger)

		user, err := u.Create(context.TODO(), mockUser)

		assert.NoError(t, err)
		assert.NotNil(t, user)

		mockUserRepository.AssertExpectations(t)
	})
}
