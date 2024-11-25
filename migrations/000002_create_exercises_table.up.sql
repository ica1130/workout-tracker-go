CREATE TABLE IF NOT EXISTS exercises (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    category text NOT NULL,
    description text
);