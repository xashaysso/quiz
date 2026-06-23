-- +goose Up
CREATE TABLE IF NOT EXISTS quiz_attempts(
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		quiz_id INT NOT NULL,
        score INT NOT NULL,
		passed_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_quiz_attempts_user_id ON quiz_attempts(user_id);
CREATE INDEX idx_quiz_attempts_quiz_id ON quiz_attempts(quiz_id);

-- +goose Down
DROP TABLE IF EXISTS quiz_attempts;
