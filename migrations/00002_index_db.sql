-- +goose Up
CREATE INDEX pack_id_removed_index ON packs (id, removed);

CREATE INDEX packs_events_index ON packs_events (pack_id);

-- +goose Down
DROP INDEX pack_id_removed_index;
DROP INDEX packs_events_index;