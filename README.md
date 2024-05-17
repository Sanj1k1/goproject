# Dota API
This is a RESTful API built in Go that allows you to manage Dota 2 databases. It uses the Gorilla Mux router and JSON for request and response bodies.

## Team Members
Amirgali Sanjar 22B030621

# API Endpoints
## Healthcheck
+ `GET /v1/healthcheck` - Checks the health of the API.
## Characters
+ `GET /v1/characters` - Retrieves characters.
+ `POST /v1/characters` - Creates character.
+ `GET /v1/characters/:id` - Retrieves a character by ID.
+ `PUT /v1/characters/:id` - Updates a character by ID.
+ `DELETE /v1/characters/:id` - DELETES a character by ID.
## Players
+ `GET /v1/players` - Retrieves players.
+ `POST /v1/players` - Creates player.
+ `GET /v1/players/:id` - Retrieves a player by ID.
+ `PUT /v1/players/:id` - Updates a player by ID.
+ `DELETE /v1/players/:id` - DELETES a player by ID.
# Database Structure 
Characters 
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
Users
```
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);
```
Players
```
CREATE TABLE IF NOT EXISTS players (
    playerid bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    nicknames text NOT NULL,
    mmr  integer NOT NULL,
    winrate  integer NOT NULL,
    totalmatches  integer NOT NULL,
    roles text[] NOT NULL
);
```
Tokens
```
CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
)
```
