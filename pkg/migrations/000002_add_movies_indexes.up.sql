CREATE INDEX IF NOT EXISTS characters_names_idx ON characters USING GIN (to_tsvector('simple', names));
CREATE INDEX IF NOT EXISTS characters_roles_idx ON characters USING GIN (roles);