package domain

import (
	"context"
	twitter "github.com/mobamoh/twitter-go-graphql"
	"github.com/mobamoh/twitter-go-graphql/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	validInput := twitter.RegisterInput{
		UserName:        "Mo",
		Email:           "mo@mail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("can register", func(t *testing.T) {
		ctx := context.Background()
		userRepo := mocks.NewUserRepo(t)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(twitter.User{ID: "666"}, nil)

		authSvc := NewAuthService(userRepo)
		res, err := authSvc.Register(ctx, validInput)
		require.NoError(t, err)
		require.NotEmpty(t, res.User.ID)
		require.NotEmpty(t, res.AccessToken)

		userRepo.AssertExpectations(t)
	})

	t.Run("username taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := mocks.NewUserRepo(t)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitter.User{}, nil)

		authSvc := NewAuthService(userRepo)
		_, err := authSvc.Register(ctx, validInput)
		require.ErrorIs(t, err, twitter.ErrUsernameTaken)

		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := mocks.NewUserRepo(t)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{}, nil)

		authSvc := NewAuthService(userRepo)
		_, err := authSvc.Register(ctx, validInput)
		require.ErrorIs(t, err, twitter.ErrEmailTaken)

		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})

	t.Run("creation failed", func(t *testing.T) {
		ctx := context.Background()
		userRepo := mocks.NewUserRepo(t)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrNotFound)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(twitter.User{}, twitter.ErrServer)

		authSvc := NewAuthService(userRepo)
		_, err := authSvc.Register(ctx, validInput)
		require.ErrorIs(t, err, twitter.ErrServer)
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		invalidInput := twitter.RegisterInput{
			UserName:        "M",
			Email:           "mo",
			Password:        "password",
			ConfirmPassword: "wrong",
		}

		ctx := context.Background()
		userRepo := mocks.NewUserRepo(t)

		authSvc := NewAuthService(userRepo)
		_, err := authSvc.Register(ctx, invalidInput)
		require.ErrorIs(t, err, twitter.ErrValidation)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})
}