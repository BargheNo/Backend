package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type City struct {
	database.Model
	Name       string
	ProvinceID uint
}
