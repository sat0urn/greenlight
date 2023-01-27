CREATE TABLE IF NOT EXISTS directors (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    surname text NOT NULL,
    awards text[] NOT NULL
);