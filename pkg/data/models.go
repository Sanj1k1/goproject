package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")

)

type Models struct {
	Characters interface {
		Insert(character *Character) error
		Get(id int64) (*Character, error)
		Update(character *Character) error
		Delete(id int64) error
		GetAll(Name string, Roles []string, filters Filters) ([]*Character, Metadata,error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Characters: MockCharacterModel{DB: db},
	}
}
