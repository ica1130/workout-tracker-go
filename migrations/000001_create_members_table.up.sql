CREATE TABLE IF NOT EXISTS members (
    id bigserial PRIMARY KEY,
    email text UNIQUE NOT NULL,
    name text NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    height int,
    weight int,
    created_at timestamp(0) with time zone DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);