package services

import (
	"errors"
	"fmt"
	"time"

	"authentication/backend/internal/models"
	"authentication/backend/internal/repositories"
	"authentication/backend/internal/utils"

	"github.com/google/uuid"
)

type AuthService struct {
	repo *repositories.UserRepo
}

func NewAuthService(repo *repositories.UserRepo) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) Register(email, password string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	return s.repo.Create(&models.User{
		Email:    email,
		Password: hash,
		Role:     "user",
	})
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByMail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// 🔐 Role-aware JWT
	return utils.GenerateToken(user.ID, user.Role)
}

func (s *AuthService) ForgotPassword(email string) (bool, error) {
	user, err := s.repo.FindByMail(email)
	if err != nil {
		return false, nil
	}

	token := uuid.NewString()
	expiry := time.Now().Add(15 * time.Minute)

	user.ResetToken = &token
	user.ResetTokenExpiry = &expiry

	if err := s.repo.Update(user); err != nil {
		return true, err
	}

	if err := sendEmail(user.Email, token); err != nil {
		fmt.Println("email error:", err)
		return true, err
	}

	return true, nil
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("old password incorrect")
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hash
	return s.repo.Update(user)
}

func (s *AuthService) GetUserByResetToken(token string) (*models.User, error) {
	user, err := s.repo.FindByResetToken(token)
	if err != nil || user == nil {
		return nil, errors.New("invalid token")
	}

	if user.ResetTokenExpiry == nil || user.ResetTokenExpiry.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return user, nil
}

func (s *AuthService) UpdatePassword(userID uint, hashedPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Password = hashedPassword
	return s.repo.Update(user)
}

func (s *AuthService) ClearResetToken(userID uint) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.ResetToken = nil
	user.ResetTokenExpiry = nil
	return s.repo.Update(user)
}

// Promote existing user to admin (manual / internal use)
func (s *AuthService) PromoteToAdmin(userID uint) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Role = "admin"
	return s.repo.Update(user)
}
