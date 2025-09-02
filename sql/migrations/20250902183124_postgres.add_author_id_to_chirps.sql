-- +goose Up
ALTER TABLE chirps
ADD COLUMN author_id UUID NOT NULL REFERENCES users(id);

-- +goose Down
ALTER TABLE chirps
DROP COLUMN author_id;