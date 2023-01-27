CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    role text NOT NULL DEFAULT 'user',
    version integer NOT NULL DEFAULT 1
);