package user

import (
	"context"
	"errors"
	"go-api/domain/entity"
	"go-api/domain/repository"
)

type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateUserUseCase(
	userRepo repository.UserRepository,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
	}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, clerkID string, firstName string, lastName string, banned bool, email string) (*entity.User, error) {
	user := entity.User{
		ClerkID:   clerkID,
		FirstName: firstName,
		LastName:  lastName,
		Banned:    banned,
		Email:     email,
	}

	err := u.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}
