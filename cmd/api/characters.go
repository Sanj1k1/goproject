package main

import (
	"errors" // New import
	"fmt"
	"goproject/pkg/data"
	"net/http"
)


func (app *application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	
	var input struct {
		Name      string   `json:"names"`
		Health    int32    `json:"health"`
		MoveSpeed int32    `json:"movespeed"`
		Mana      int32    `json:"mana"`
		Roles     []string `json:"roles"`
	}
	
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character := &data.Character{
		Name:      input.Name,
		Health:    input.Health,
		MoveSpeed: input.MoveSpeed,
		Mana:      input.Mana,
		Roles:     input.Roles,
	}

	err = app.models.Characters.Insert(character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}


	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/characters/%d", character.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"character": character}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}


func (app *application) showCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.Get(id)
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
		Name      string   `json:"names"`
		Health    int32    `json:"health"`
		MoveSpeed int32    `json:"movespeed"`
		Mana      int32    `json:"mana"`
		Roles     []string `json:"roles"`
	}
	
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character.Name = input.Name
	character.Health = input.Health
	character.MoveSpeed = input.MoveSpeed
	character.Mana = input.Mana
	character.Roles = input.Roles

	err = app.models.Characters.Update(character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Characters.Delete(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "character successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	}
	