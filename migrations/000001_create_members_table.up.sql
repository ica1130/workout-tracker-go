CREATE TABLE IF NOT EXISTS members (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    height int,
    weight int,
    created_at timestamp(0) with time zone DEFAULT NOW()
);