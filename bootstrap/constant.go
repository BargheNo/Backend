package bootstrap

import "fmt"

type Constants struct {
	Context      Context
	LogLevel     LogLevel
	RedisKey     RedisKey
	Field        ErrorField
	Tag          ErrorTag
	SMSTemplates SMSTemplates
	JWTKeysPath  JWTKeysPath
}

type Context struct {
	Translator                   string
	IsLoadedValidationTranslator string
}

type LogLevel struct {
	Debug string
	Info  string
	Warn  string
	Error string
	Fatal string
}

type RedisKey struct {
}

type ErrorField struct {
	Phone    string
	Password string
	OTP      string
}

type ErrorTag struct {
	AlreadyRegistered      string
	MinimumLength          string
	ContainsLowercase      string
	ContainsUppercase      string
	ContainsNumber         string
	ContainsSpecialChar    string
	OTPExpired             string
	InvalidOTP             string
	NotRegistered          string
	InvalidAuthCredentials string
	ExpiredAuthToken       string
	InvalidAuthToken       string
	Unauthorized           string
}

type SMSTemplates struct {
	OTP string
}

type JWTKeysPath struct {
	PublicKey  string
	PrivateKey string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator:                   "translator",
			IsLoadedValidationTranslator: "isLoadedValidationTranslator",
		},
		LogLevel: LogLevel{
			Debug: "debug",
			Info:  "info",
			Warn:  "warn",
			Error: "error",
			Fatal: "fatal",
		},
		Field: ErrorField{
			Phone:    "phone",
			Password: "password",
			OTP:      "otp",
		},
		Tag: ErrorTag{
			AlreadyRegistered:      "alreadyRegistered",
			MinimumLength:          "minimumLength",
			ContainsLowercase:      "containsLowercase",
			ContainsUppercase:      "containsUppercase",
			ContainsNumber:         "containsNumber",
			ContainsSpecialChar:    "containsSpecialChar",
			OTPExpired:             "otpExpired",
			InvalidOTP:             "invalidOTP",
			NotRegistered:          "notRegistered",
			InvalidAuthCredentials: "invalidAuthCredentials",
			ExpiredAuthToken:       "expiredAuthToken",
			InvalidAuthToken:       "invalidAuthToken",
			Unauthorized:           "unauthorized",
		},
		SMSTemplates: SMSTemplates{
			OTP: "sendOTPTemplate",
		},
		JWTKeysPath: JWTKeysPath{
			PublicKey:  "./internal/application/adapter/jwt/publicKey.pem",
			PrivateKey: "./internal/application/adapter/jwt/privateKey.pem",
		},
	}
}

func (r *RedisKey) GenerateOTPKey(value string) string {
	return fmt.Sprintf("otp:%s", value)
}
