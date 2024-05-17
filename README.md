# Dota API
The Dota 2 API provides developers with access to data and functionalities of the Dota 2 game, enabling the creation of applications and services related to the Dota 2 universe. The API offers endpoints for retrieving information about characters, players, as well as the ability to create and update data. With this API, developers can integrate character and player information into their applications, as well as automate data management processes within the game.
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
