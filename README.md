# Dota API
This is a RESTful API built in Go that allows you to manage Dota 2 databases. It uses the Gorilla Mux router and JSON for request and response bodies.

# API Endpoints
## Healthcheck
+ `GET /v1/health` - Checks the health of the API.
## Characters
+ `POST /v1/characters` - Creates character.
+ `GET /v1/characters/:id` - Retrieves a character by ID.
+ `PUT /v1/characters/:id` - Updates a character by ID.
+ `DELETE /v1/characters/:id` - DELETES a character by ID.

# Database Structure 
```
CREATE TABLE IF NOT EXISTS characters (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    names text NOT NULL,
    health  integer NOT NULL,
    movespeed  integer NOT NULL,
    mana  integer NOT NULL,
    roles text[] NOT NULL
);
```
## Amirgali Sanjar 22B030621
