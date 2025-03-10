package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type SampleRepository interface {
	Create(db database.Database, user *entity.User) error
	Delete(db database.Database, userID uint) error
}
