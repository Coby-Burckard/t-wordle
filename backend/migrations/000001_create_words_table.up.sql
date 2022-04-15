CREATE TABLE IF NOT EXISTS words (
  id bigserial PRIMARY KEY,
  answer text NOT NULL,
  solve_count integer NOT NULL DEFAULT 0,
  solve_time integer NOT NULL DEFAULT 0
);

ALTER TABLE words ADD CONSTRAINT solve_count_check CHECK (solve_count >= 0);

ALTER TABLE words ADD CONSTRAINT solve_time_check CHECK (solve_time >= 0);

ALTER TABLE words ADD CONSTRAINT answer_length_check CHECK (char_length(answer) = 5);