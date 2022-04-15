CREATE TABLE IF NOT EXISTS guesses(
  ID bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  word_id int NOT NULL,
  attempt_id int NOT NULL,
  submission text NOT NULL,
  submission_time int NOT NULL
);

-- TODO: enforce word_id + user uniqueness

ALTER TABLE guesses ADD CONSTRAINT fk_word FOREIGN KEY (word_id) REFERENCES words(id);

ALTER TABLE guesses ADD CONSTRAINT fk_attempt FOREIGN KEY (attempt_id) REFERENCES attempts(id);

ALTER TABLE guesses ADD CONSTRAINT submission_length_check CHECK (char_length(submission) = 5);