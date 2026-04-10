package api

type AuthorizeResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type PublicKeyResponse struct {
	Key string `json:"key"`
}
