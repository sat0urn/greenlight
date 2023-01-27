CREATE INDEX IF NOT EXISTS directors_name_idx ON directors USING gin (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS directors_awards_idx ON directors USING gin (awards);