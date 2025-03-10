package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type User struct {
	database.Model
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	Password     string
	NationalCode string
}
