CREATE TABLE IF NOT EXISTS characters (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    names text NOT NULL,
    health  integer NOT NULL,
    movespeed  integer NOT NULL,
    mana  integer NOT NULL,
    roles text[] NOT NULL
);