package data

import (
	"fmt"
	"context" 
	"database/sql"
	"errors"
	"time"
	"goproject/pkg/validator"
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

func ValidateCharacter(v *validator.Validator, character *Character) {
	v.Check(character.Name != "", "Name", "must be provided")
	v.Check(len(character.Name) <= 500, "Name", "must not be more than 500 bytes long")
	v.Check(character.Health != 0, "Health", "must be provided")
	v.Check(character.Health > 0, "Health", "must be greater than 0")
	v.Check(character.MoveSpeed != 0, "MoveSpeed", "must be provided")
	v.Check(character.MoveSpeed > 0, "MoveSpeed", "must be a positive integer")
	v.Check(character.Roles != nil, "Roles", "must be provided")
	v.Check(len(character.Roles) >= 1, "Roles", "must contain at least 1 genre")
	v.Check(validator.Unique(character.Roles), "genres", "must not contain duplicate values")
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.DB.QueryRowContext(ctx,query, args...).Scan(&character.ID, &character.CreatedAt)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := c.DB.QueryRowContext(ctx,query, id).Scan(
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
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx,query, args...).Scan(&character.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
		return ErrEditConflict
		default:
		return err
		}
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx,query, id)
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

func (c MockCharacterModel) GetAll(Name string, Roles []string, filters Filters) ([]*Character,Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(),id, created_at, names, health, movespeed, mana,roles
		FROM characters
		WHERE (to_tsvector('simple', names) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (roles @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{Name, pq.Array(Roles), filters.limit(), filters.offset()}

	rows, err := c.DB.QueryContext(ctx, query,args...)
	if err != nil {
		return nil,Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	characters := []*Character{}

	for rows.Next() {

		var character Character
		
		err := rows.Scan(
			&totalRecords,
			&character.ID,
			&character.CreatedAt,
			&character.Name,
			&character.Health,
			&character.MoveSpeed,
			&character.Mana,
			pq.Array(&character.Roles),
		)
		if err != nil {
			return nil, Metadata{},err
		}
		
		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{} ,err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return characters, metadata,nil
}
