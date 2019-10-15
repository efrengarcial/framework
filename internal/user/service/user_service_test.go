package service

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/efrengarcial/framework/internal/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger = log.New()


func init() {
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	logger.Formatter = &log.JSONFormatter{}
}

//https://github.com/bxcodec/go-clean-arch/blob/master/article/usecase/article_ucase_test.go
// mockery -name=Repository
func TestInsert(t *testing.T) {
	m :=  new(mocks.Repository)
	var mockUser domain.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)
	mockUser.ID = 0

	t.Run("success", func(t *testing.T) {
		m.On("Insert", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()
		m.On("FindOneByLogin", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		m.On("FindOneByEmail", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		u := NewService(m,  logger)
		err := u.Create(context.TODO(), &mockUser)
		assert.NoError(t, err)
		m.AssertExpectations(t)
	})
}
