package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"tracker-server/internal/pkg/uuid"
	"tracker-server/internal/repo"
)

type UserService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(ctx context.Context, email, password, fullName, displayName string) (*repo.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	user := &repo.User{
		ID:          id,
		Email:       email,
		FullName:    fullName,
		DisplayName: displayName,
		CreatedAt:   time.Now(),
	}

	secret := &repo.UserSecret{
		ID:    id,
		Value: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user, secret); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*repo.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
