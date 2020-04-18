package domain

// LoginRequest represent user login payload
type LoginRequest struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	AccessTokenExpired int    `json:"accessTokenExpired"`
}

// LoginResponse represent user login response
type LoginResponse struct {
	Email              string `json:"email"`
	AccessToken        string `json:"accessToken"`
	AccessTokenExpired int64  `json:"accessTokenExpired"`
}
