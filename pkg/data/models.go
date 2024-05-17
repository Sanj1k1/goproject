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
	Players interface{
		Insert(player *Player) error
		Get(playerid int64) (*Player, error)
		Update(player *Player) error
		Delete(playerid int64) error
		GetAll(Nickname string, Roles []string, filters Filters) ([]*Player,Metadata,error)
	}
	Users UserModel 
	Tokens TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Characters: MockCharacterModel{DB: db},
		Players: MockPlayerModel{DB: db},
		Permissions: PermissionModel{DB: db}, 
		Tokens: TokenModel{DB: db},
		Users: UserModel{DB: db},
	}
}