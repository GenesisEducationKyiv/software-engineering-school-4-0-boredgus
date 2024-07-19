package db

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type database struct {
	db *gorm.DB
}

func NewDatabase(url, schemaName string) (*database, error) {
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

	return &database{db: db}, nil
}

func (d *database) Create(value interface{}) error {
	return d.db.Create(value).Error
}

func (d *database) Delete(value interface{}) error {
	return d.db.Delete(value).Error
}

func (d *database) Where(query any, args ...any) repo.DB {
	return &database{db: d.db.Where(query, args...)}
}
