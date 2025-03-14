package bootstrap

import "fmt"

type Constants struct {
	Context      Context
	LogLevel     LogLevel
	RedisKey     RedisKey
	Field        ErrorField
	Tag          ErrorTag
	SMSTemplates SMSTemplates
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
}

type ErrorTag struct {
	AlreadyRegistered   string
	MinimumLength       string
	ContainsLowercase   string
	ContainsUppercase   string
	ContainsNumber      string
	ContainsSpecialChar string
}

type SMSTemplates struct {
	OTP string
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
		},
		Tag: ErrorTag{
			AlreadyRegistered:   "alreadyRegistered",
			MinimumLength:       "minimumLength",
			ContainsLowercase:   "containsLowercase",
			ContainsUppercase:   "containsUppercase",
			ContainsNumber:      "containsNumber",
			ContainsSpecialChar: "containsSpecialChar",
		},
		SMSTemplates: SMSTemplates{
			OTP: "sendOTPTemplate",
		},
	}
}

func (r *RedisKey) GenerateOTPKey(value string) string {
	return fmt.Sprintf("otp:%s", value)
}
