package data

import (
	"context"
	"database/sql"
	"errors"
	"goproject/pkg/validator"
	"time"

	"github.com/lib/pq"
)

type Player struct {
	PlayerID 	int64	`json:"playerid"`
	CreatedAt 	time.Time	`json:"created_at"`
	Nickname 	string	`json:"nicknames"`
	MMR 	int32	`json:"mmr"`
	WinRate  int64	`json:"winrate"`
	TotalMatches	int64	`json:"totalmatches"`
	Roles 	[]string	`json:"roles"`
}

func ValidatePlayer(v *validator.Validator, player *Player) {
	v.Check(player.Nickname != "", "Nickname", "must be provided")
	v.Check(len(player.Nickname) <= 500, "Nickname", "must not be more than 500 bytes long")

	v.Check(player.MMR != 0, "MMR", "must be provided")
	v.Check(player.MMR > 0, "MMR", "must be greater than 0")

	v.Check(player.WinRate != 0, "WinRate", "must be provided")
	v.Check(player.WinRate > 0, "WinRate", "must be a positive integer")

	v.Check(player.Roles != nil, "Roles", "must be provided")
	v.Check(len(player.Roles) >= 1, "Roles", "must contain at least 1 genre")
	v.Check(validator.Unique(player.Roles), "Roles", "must not contain duplicate values")
}

type MockPlayerModel struct {
	DB *sql.DB
}

func (p MockPlayerModel) Insert(player *Player) error {
	query := `
			INSERT INTO players (nicknames, mmr, winrate, totalmatches,roles)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING playerid, created_at`
	
	args := []interface{}{player.Nickname, player.MMR, player.WinRate,player.TotalMatches, pq.Array(player.Roles)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return p.DB.QueryRowContext(ctx,query, args...).Scan(&player.PlayerID, &player.CreatedAt)
}

func (p MockPlayerModel) Get(playerid int64) (*Player, error) {
	if playerid < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT playerid, created_at, nicknames, mmr, winrate, totalmatches ,roles
		FROM players
		WHERE playerid = $1`

	var player Player

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := p.DB.QueryRowContext(ctx,query, playerid).Scan(
		&player.PlayerID,
		&player.CreatedAt,
		&player.Nickname,
		&player.MMR,
		&player.WinRate,
		&player.TotalMatches,
		pq.Array(&player.Roles),
	)
		
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &player, nil
}

func (p MockPlayerModel) Update(player *Player) error {
	query := `
	UPDATE players
	SET nicknames = $1, mmr = $2, winrate = $3, totalmatches = $4, roles=$5
	WHERE playerid = $6
	RETURNING playerid, created_at, nicknames, mmr, winrate, totalmatches, roles`

	args := []interface{}{
		player.Nickname,
		player.MMR,
		player.WinRate,
		player.TotalMatches,
		pq.Array(player.Roles),
		player.PlayerID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx,query, args...).Scan(
		&player.PlayerID,
		&player.CreatedAt,
		&player.Nickname,
		&player.MMR,
		&player.WinRate,
		&player.TotalMatches,
		pq.Array(&player.Roles),
	)
	
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return ErrEditConflict
			default:
				return err
		}
	}
	return nil
}

func (p MockPlayerModel) Delete(playerid int64) error {
	if playerid < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM players
	WHERE playerid = $1`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx,query, playerid)
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