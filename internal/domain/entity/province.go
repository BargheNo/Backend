package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Province struct {
	database.Model
	Name   string
	Cities []City `gorm:"foreignKey:ProvinceID"`
}
