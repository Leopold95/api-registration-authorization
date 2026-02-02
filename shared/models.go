package shared

import "github.com/google/uuid"

type UserModel struct {
	Id       uuid.UUID
	Email    string
	Name     string
	Password string
}

type TokenModel struct {
	Token        string
	RefreshToken string
}
