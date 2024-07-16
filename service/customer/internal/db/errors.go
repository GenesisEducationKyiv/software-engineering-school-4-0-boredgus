package db

import (
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/repo"
	"gorm.io/gorm"
)

var errorsMap = map[repo.DBError]error{
	repo.UniqueViolationErr: gorm.ErrDuplicatedKey,
}

func IsError(err error, targetErr repo.DBError) bool {
	return errors.Is(err, errorsMap[targetErr])
}
