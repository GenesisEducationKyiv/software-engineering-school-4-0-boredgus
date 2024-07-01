package err

import "errors"

var (
	InvalidArgumentErr = errors.New("invalid argument")
	NotFoundErr        = errors.New("not found")
	UniqueViolationErr = errors.New("unique violation")
)
