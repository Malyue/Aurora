package user

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type RegisterResponse struct {
}

type GetUserInfoRequest struct {
	ID string
}

type GetUserInfoResponse struct {
}
