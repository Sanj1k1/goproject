package main
import (
	"fmt"
	"net/http"
	"goproject/pkg/data"
	"goproject/pkg/validator"
	"errors" 
)

func (app *application) createPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nickname 	string	`json:"nicknames"`
		MMR 	int32	`json:"mmr"`
		WinRate  int64	`json:"winrate"`
		TotalMatches	int64	`json:"totalmatches"`
		Roles     []string  `json:"roles"`	
	}
	
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	player := &data.Player{
		Nickname: 	input.Nickname,
		MMR:	 input.MMR,
		WinRate: 	input.WinRate,
		TotalMatches:	 input.TotalMatches,
		Roles: 	input.Roles,
	}

	v := validator.New()

	if data.ValidatePlayer(v, player); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}


	err = app.models.Players.Insert(player)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/players/%d", player.PlayerID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"player": player}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPlayerHandler(w http.ResponseWriter, r *http.Request) {

	playerid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}		

	player, err := app.models.Players.Get(playerid)
	if err != nil {
	switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	playerid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	player, err := app.models.Players.Get(playerid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Nickname     *string  `json:"nicknames"`
		MMR          *int32   `json:"mmr"`
		WinRate      *int64   `json:"winrate"`
		TotalMatches *int64   `json:"totalmatches"`
		Roles        []string `json:"roles"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Nickname != nil {
		player.Nickname = *input.Nickname 
	}
	if input.MMR != nil {
		player.MMR = *input.MMR 
	}
	if input.WinRate != nil {
		player.WinRate = *input.WinRate 
	}
	if input.TotalMatches != nil {
		player.TotalMatches = *input.TotalMatches 
	}
	if input.Roles != nil {
		player.Roles = input.Roles 
	}

	v := validator.New()
	if data.ValidatePlayer(v, player); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Players.Update(player)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrEditConflict):
				app.editConflictResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	playerid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Players.Delete(playerid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "player successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	}
	