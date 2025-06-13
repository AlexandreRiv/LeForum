package service

import (
	"LeForum/internal/domain"
	"LeForum/internal/storage/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SaveUserIfNotExists(email, username string) error {
	return s.repo.SaveUserIfNotExists(email, username)
}

func (s *UserService) CreateUser(username, email, password string) error {
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(username, email, hashedPassword)
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *UserService) GetUserStats(email string) (postCount, responseCount, likeCount int, err error) {
	return s.repo.GetUserStats(email)
}

func (s *UserService) UpdateUserRole(userID int, role domain.RoleType) error {
	return s.repo.UpdateUserRole(userID, string(role))
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

// Complétons le système dans GetCurrentUser pour inclure le rôle dans la session
func (s *UserService) GetUserWithRole(email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(email)
}
