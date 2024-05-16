package main

import (
	"errors"
	"fmt"
	"goproject/pkg/data"
	"goproject/pkg/validator"
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

	v := validator.New()

	if data.ValidateCharacter(v, character); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
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
		Name      *string   `json:"names"`
		Health    *int32    `json:"health"`
		MoveSpeed *int32    `json:"movespeed"`
		Mana      *int32    `json:"mana"`
		Roles     []string `json:"roles"`
	}
	
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		character.Name = *input.Name
	}

	if input.Health != nil {
		character.Health = *input.Health
	}

	if input.MoveSpeed != nil {
		character.MoveSpeed = *input.MoveSpeed
	}

	if input.Mana != nil {
		character.Mana = *input.Mana
	}

	if input.Roles != nil {
		character.Roles = input.Roles
	}
	

	v := validator.New()
	if data.ValidateCharacter(v, character); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Characters.Update(character)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
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

func (app *application) listCharactersHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string
		Roles []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	
	input.Name = app.readString(qs, "name", "")
	input.Roles = app.readCSV(qs, "roles", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "names", "health", "mana","movespeed", "-id","roles", "-names", "-health", "-mana","-movespeed","-roles"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	characters,metadata, err := app.models.Characters.GetAll(input.Name, input.Roles, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"characters": characters,"metadata":metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
	