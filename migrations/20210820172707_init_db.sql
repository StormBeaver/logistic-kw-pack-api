-- +goose Up
CREATE TABLE pack (
  id BIGSERIAL PRIMARY KEY,
  foo BIGINT NOT NULL
);

-- +goose Down
DROP TABLE pack;
