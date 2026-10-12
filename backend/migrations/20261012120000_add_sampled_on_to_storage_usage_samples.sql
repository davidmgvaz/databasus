-- +goose Up
-- +goose StatementBegin
-- Installations that already ran 20260928110000 have no sampled_on column. The sampler now
-- keeps one sample per storage per UTC day, so backfill the day from sampled_at, keep the
-- latest reading of each day, and only then enforce uniqueness.
ALTER TABLE storage_usage_samples
    ADD COLUMN sampled_on DATE;

UPDATE storage_usage_samples
SET sampled_on = (sampled_at AT TIME ZONE 'UTC')::date;

DELETE FROM storage_usage_samples AS older
USING storage_usage_samples AS newer
WHERE older.storage_id = newer.storage_id
  AND older.sampled_on = newer.sampled_on
  AND (older.sampled_at, older.id) < (newer.sampled_at, newer.id);

ALTER TABLE storage_usage_samples
    ALTER COLUMN sampled_on SET NOT NULL;

ALTER TABLE storage_usage_samples
    ADD CONSTRAINT uq_storage_usage_samples_storage_id_sampled_on
    UNIQUE (storage_id, sampled_on);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE storage_usage_samples
    DROP CONSTRAINT IF EXISTS uq_storage_usage_samples_storage_id_sampled_on;

ALTER TABLE storage_usage_samples
    DROP COLUMN IF EXISTS sampled_on;
-- +goose StatementEnd
