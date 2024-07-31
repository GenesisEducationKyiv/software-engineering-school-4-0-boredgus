package db

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDatabase(url, schemaName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s.", schemaName),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to postgres db: %w", err)
	}

	if err = db.AutoMigrate(&repo.Customer{}); err != nil {
		return nil, fmt.Errorf("failed to automigrate model 'customer': %w", err)
	}

	return db, nil
}
