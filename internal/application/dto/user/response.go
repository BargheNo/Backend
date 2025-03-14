package userdto

type OTPData struct {
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type UserInfoResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Permissions  []string `json:"permissions"`
}
