package seed

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ContactTypeSeeder struct {
	corporationRepository repository.CorporationRepository
	db                    database.Database
}

func NewContactTypeSeeder(
	corporationRepository repository.CorporationRepository,
	db database.Database,
) *ContactTypeSeeder {
	return &ContactTypeSeeder{
		corporationRepository: corporationRepository,
		db:                    db,
	}
}

var contactTypes = []string{
	"شماره تلفن",
	"ایمیل",
	"ایتا",
	"بله",
	"تارنما",
	"WhatsApp",
	"Instagram",
	"LinkedIn",
	"Telegram",
}

func (seeder *ContactTypeSeeder) SeedContactTypes() {
	for _, name := range contactTypes {
		_, exist := seeder.corporationRepository.FindContactTypeByName(seeder.db, name)
		if exist {
			continue
		}
		contactType := &entity.ContactType{
			Name: name,
		}
		if err := seeder.corporationRepository.CreateContactType(seeder.db, contactType); err != nil {
			panic(err)
		}
	}
}
