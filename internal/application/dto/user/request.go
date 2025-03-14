package userdto

type BasicRegisterRequest struct {
	FirstName string
	LastName  string
	Phone     string
	Password  string
}

type VerifyPhoneRequest struct {
	Phone string
	OTP   string
}

type LoginRequest struct {
	Phone    string
	Password string
}
