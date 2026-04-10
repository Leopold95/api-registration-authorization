package domain

type UserModel struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	ProfileId string `json:"profileId"`
	Id        string `json:"id"`
}
