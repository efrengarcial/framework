package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/mocks"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGet(t *testing.T) {
	mockRepository :=  new(mocks.Repository)

	mockUser := &model.User {
		FirstName: "Efren",
		LastName: "Garcia",
		Email: "efren.gl@gmail.com",
		Login: "efren",
	}

	t.Run("success", func(t *testing.T) {
		mockRepository.On("Get", mock.AnythingOfType("int64")).Return(mockUser, nil).Once()
		tokenService := service.NewTokenService()
		u := service.NewUserService(mockRepository ,tokenService)

		user, err := u.Get(context.TODO(), mockUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, user)

		mockRepository.AssertExpectations(t)
	})
}