CREATE TABLE IF NOT EXISTS workout_details (
    id SERIAL PRIMARY KEY,
    workout_id bigint NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id bigint NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    set int NOT NULL,
    repetitions int NOT NULL,
    weight float NOT NULL
);