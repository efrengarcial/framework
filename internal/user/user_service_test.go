package user

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"os"
	"testing"

	"github.com/efrengarcial/framework/internal/user/mocks"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	mockAuthority := domain.Authority{Model: domain.Model{ID: 1}}

	mockUser := &domain.User{
		FirstName: "Juan",
		LastName: "Perez",
		Email: "juan.perez@gmail.com",
		Login: "jperez",
		Authorities: []domain.Authority{mockAuthority },
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
