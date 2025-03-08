package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Server    Server
	Logger    Logger
	RateLimit RateLimit
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

func NewEnvironments() *Env {
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
	}
}
