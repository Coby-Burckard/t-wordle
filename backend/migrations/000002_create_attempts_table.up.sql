CREATE TABLE IF NOT EXISTS attempts (
  id bigserial PRIMARY KEY,
  word_id integer NOT NULL,
  is_open boolean NOT NULL DEFAULT TRUE,
  is_solved boolean NOT NULL DEFAULT FALSE,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

-- TODO: enforce word_id + user uniqueness when user added

ALTER TABLE attempts ADD CONSTRAINT fk_word FOREIGN KEY(word_id) REFERENCES words(id);