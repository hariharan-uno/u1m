
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE ranking (
  name VARCHAR(255) NOT NULL,
  rank INTEGER NOT NULL,
  day DATE NOT NULL,
  UNIQUE rank_day (rank, day),
  INDEX (name)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE ranking;
