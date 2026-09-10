package application

import "golang.org/x/crypto/bcrypt"

type PasswordHasherService struct {
}

func NewPasswordHasherService() *PasswordHasherService {
	return &PasswordHasherService{}
}

func (this *PasswordHasherService) HashPassword(rawPassword string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

func (this *PasswordHasherService) VerifyPassword(rawPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	return err == nil
}
