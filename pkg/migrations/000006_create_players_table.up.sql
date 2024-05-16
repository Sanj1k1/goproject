CREATE TABLE IF NOT EXISTS players (
    playerid bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    nicknames text NOT NULL,
    mmr  integer NOT NULL,
    winrate  integer NOT NULL,
    totalmatches  integer NOT NULL,
    roles text[] NOT NULL
);