CREATE TABLE IF NOT EXISTS workouts (
    id bigserial PRIMARY KEY,
    member_id bigint NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);