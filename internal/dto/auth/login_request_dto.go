package auth_dto

type LoginRQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRS struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRQ struct {
	RefreshToken string `json:"refreshToken"`
}

type CreateAccountRQ struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ForgotPasswordRQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
