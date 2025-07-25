package infrastructure

import (
	domain "task_manager/Domain"
	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
}
func NewPasswordService() domain.IPasswordService {
	return &PasswordService{}
}

func (ps *PasswordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	password = string(hashedPassword)
	return password, nil
}

func (ps *PasswordService) VerifyPassword(user *domain.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

