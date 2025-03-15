package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Server       Server
	Logger       Logger
	RateLimit    RateLimit
	PrimaryDB    Database
	PrimaryRedis Redis
	OTP          OTP
	SMSGateway   SMSGateway
}

type Server struct {
	Port string
	Mode string
}

type Logger struct {
	LogLevel      string
	LogFile       string
	ConsoleOutput string
}

type RateLimit struct {
	Limit string
	Burst string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Redis struct {
	Port      string
	Address   string
	Password  string
	RDBNumber string
}

type OTP struct {
	Length       string
	ExpiryMinute string
	MaxAttempts  string
}

type SMSGateway struct {
	APIKey string
}

func NewEnvironments() *Env {
	// godotenv.Load("../../.env")
	godotenv.Load(".env")
	return &Env{
		Server: Server{
			Port: os.Getenv("SERVER_PORT"),
			Mode: os.Getenv("SERVER_MODE"),
		},
		Logger: Logger{
			LogLevel:      os.Getenv("LOG_LEVEL"),
			LogFile:       os.Getenv("LOG_FILE"),
			ConsoleOutput: os.Getenv("CONSOLE_OUTPUT"),
		},
		RateLimit: RateLimit{
			Limit: os.Getenv("RATE_LIMIT"),
			Burst: os.Getenv("RATE_lIMIT_BURST"),
		},
		PrimaryDB: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		PrimaryRedis: Redis{
			Port:      os.Getenv("RDB_PORT"),
			Address:   os.Getenv("RDB_ADDRESS"),
			Password:  os.Getenv("RDB_PASSWORD"),
			RDBNumber: os.Getenv("RDB_NUMBER"),
		},
		OTP: OTP{
			Length:       os.Getenv("OTP_LENGTH"),
			ExpiryMinute: os.Getenv("OTP_EXPIRY_MINUTES"),
			MaxAttempts:  os.Getenv("OTP_MAX_ATTEMPTS"),
		},
		SMSGateway: SMSGateway{
			APIKey: os.Getenv("SMS_GATEWAY_API_KEY"),
		},
	}
}
