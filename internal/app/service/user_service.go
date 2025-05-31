package service

import (
	"errors"

	"github.com/fahrikurniawan99/studentsync-api/internal/app/repository"
	"github.com/fahrikurniawan99/studentsync-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(name, email, password string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	LoginUser(email, password string) (*model.User, error)
	GetAllUsers() ([]model.User, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) RegisterUser(name, email, password string) (*model.User, error) {
	// Periksa apakah email sudah terdaftar
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userServiceImpl) LoginUser(email, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("email atau password salah")
	}
	return user, nil
}

func (s *userServiceImpl) GetAllUsers() ([]model.User, error) {
	return s.userRepo.GetAllUsers()
}
