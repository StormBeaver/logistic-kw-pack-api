-- +goose Up
ALTER TABLE packs ADD describe VARCHAR(511) DEFAULT '';

-- +goose Down
ALTER TABLE packs
DROP describe;