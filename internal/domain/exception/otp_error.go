package exception

import "errors"

var (
	ErrOTPExpired = errors.New("OTP_EXPIRED")
	ErrInvalidOTP = errors.New("INVALID_OTP")
)

func IsOTPExpired(err error) bool {
	return errors.Is(err, ErrOTPExpired)
}

func IsInvalidOTP(err error) bool {
	return errors.Is(err, ErrInvalidOTP)
}
