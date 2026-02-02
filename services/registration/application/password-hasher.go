package application

import "golang.org/x/crypto/bcrypt"

type PasswordHasherService struct {
}

func NewPasswordHasherService() *PasswordHasherService {
	return &PasswordHasherService{}
}

func (this *PasswordHasherService) HashPassword(rawPassword string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.MinCost)
	return string(hashBytes)
}

func (this *PasswordHasherService) VerifyPassword(rawPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	return err == nil
}
