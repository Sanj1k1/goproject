package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Character struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"names"`
	Health    int32     `json:"health"`
	MoveSpeed int32     `json:"movespeed"`
	Mana      int32     `json:"mana"`
	Roles     []string  `json:"roles"`
}

type MockCharacterModel struct {
	DB *sql.DB
}


func (c MockCharacterModel) Insert(character *Character) error {
	query := `
			INSERT INTO characters (names,health,movespeed,mana,roles)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

	args := []interface{}{character.Name, character.Health, character.MoveSpeed, character.Mana, pq.Array(character.Roles)}

	return c.DB.QueryRow(query, args...).Scan(&character.ID, &character.CreatedAt)
}


func (c MockCharacterModel) Get(id int64) (*Character, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, created_at, names,health,movespeed,mana,roles
	FROM characters
	WHERE id = $1`
	var character Character

	err := c.DB.QueryRow(query, id).Scan(
		&character.ID,
		&character.CreatedAt,
		&character.Name,
		&character.Health,
		&character.MoveSpeed,
		&character.Mana,
		pq.Array(&character.Roles),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &character, nil

}

func (c MockCharacterModel) Update(character *Character) error {
	query := `
	UPDATE characters
	SET names = $1, health = $2, movespeed = $3, mana = $4,roles=$5
	WHERE id = $6
	RETURNING *`
	args := []interface{}{
		character.Name,
		character.Health,
		character.MoveSpeed,
		character.Mana,
		pq.Array(character.Roles),
		character.ID,
	}

	return c.DB.QueryRow(query, args...).Scan(&character.Roles)

}

func (c MockCharacterModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	
	query := `
	DELETE FROM characters
	WHERE id = $1`

	result, err := c.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
