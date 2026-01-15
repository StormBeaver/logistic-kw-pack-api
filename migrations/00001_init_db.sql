-- +goose Up
CREATE TABLE IF NOT EXISTS packs(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  describe VARCHAR(511) DEFAULT '',
  removed BOOLEAN DEFAULT FALSE,
  created TIMESTAMP,
  updated TIMESTAMP
);

CREATE TABLE IF NOT EXISTS packs_events(
  id BIGSERIAL PRIMARY KEY,
  pack_id BIGINT,
  type TEXT,
  lock BOOLEAN DEFAULT FALSE,
  payload JSONB,
  updated TIMESTAMP,
  FOREIGN KEY (pack_id) REFERENCES packs(id)
);

CREATE INDEX ON packs_events (pack_id)

CREATE INDEX ON packs_events (id)

-- +goose Down
DROP TABLE packs_events;
DROP TABLE packs;
