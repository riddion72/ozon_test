package dbErrors

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const (
	ErrNotFound         = "entity not found"
	ErrUniqAlreadyExist = "uniq already exists"
	ErrAlreadyExists    = "entity already exists"
	ErrFailedConnection = "failed connection"
)

func PrepareError(err error) error {
	pErr, ok := err.(*pq.Error)
	if !ok {
		switch err {
		case sql.ErrNoRows:
			return errors.New(ErrNotFound)
		case sql.ErrConnDone:
			return errors.New(ErrFailedConnection)
		default:
			return errors.New(ErrNotFound)
		}
	}

	switch pErr.Code {
	case "23503":
		return errors.New(ErrNotFound)
	case "23505":
		return errors.New(ErrUniqAlreadyExist)
	case "08006":
		return errors.New(ErrFailedConnection)
	}

	return err
}
