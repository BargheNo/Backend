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

type CredentialResponse struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	NationalID string `json:"nationalID"`
	ProfilePic string `json:"profilePic"`
	Status     string `json:"status"`
}

type UserResponse struct {
	ID uint `json:"id"`
}
