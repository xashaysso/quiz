-- +goose Up
CREATE TABLE IF NOT EXISTS quiz_global_stats(
    quiz_id INT PRIMARY KEY,
    total_attempts INT NOT NULL DEFAULT 0,
    accumulated_score BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS quiz_global_stats;
