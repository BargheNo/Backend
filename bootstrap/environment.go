package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Server Server
}

type Server struct {
	Port string
	Mode string
}

func NewEnvironments() *Env {
	godotenv.Load(".env")
	return &Env{
		Server: Server{
			Port: os.Getenv("SERVER_PORT"),
			Mode: os.Getenv("SERVER_MODE"),
		},
	}
}
