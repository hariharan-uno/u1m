
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE current (
  rank INTEGER UNIQUE NOT NULL,
  domain VARCHAR(255) NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE current;
