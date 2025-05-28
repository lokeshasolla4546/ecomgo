package services

import (
	"bego/config"
	"bego/repositories"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		Repo: repositories.NewUserRepository(config.DB),
	}
}

func (s *UserService) Register(userID, password string) (*repositories.User, error) {
	return s.Repo.Register(userID, password, "user")
}

func (s *UserService) Login(userID, password string) (*repositories.User, string, error) {
	user, err := s.Repo.Login(userID, password)
	if err != nil {
		return nil, "", err
	}

	token, err := generateJWT(user.UserID, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return user, token, nil
}

func generateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecretKey)
}
