-- +goose Up
-- +goose StatementBegin
-- One row per storage per day, read by the dashboard to fit a line through the used share
-- and project when the storage fills up.
CREATE TABLE storage_usage_samples (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    storage_id  UUID        NOT NULL,
    sampled_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    sampled_on  DATE        NOT NULL,
    total_bytes BIGINT      NOT NULL,
    used_bytes  BIGINT      NOT NULL,
    free_bytes  BIGINT      NOT NULL
);

ALTER TABLE storage_usage_samples
    ADD CONSTRAINT fk_storage_usage_samples_storage_id
        FOREIGN KEY (storage_id) REFERENCES storages (id) ON DELETE CASCADE;

ALTER TABLE storage_usage_samples
    ADD CONSTRAINT uq_storage_usage_samples_storage_id_sampled_on UNIQUE (storage_id, sampled_on);

CREATE INDEX idx_storage_usage_samples_storage_id_sampled_at
    ON storage_usage_samples (storage_id, sampled_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_storage_usage_samples_storage_id_sampled_at;

ALTER TABLE storage_usage_samples
    DROP CONSTRAINT IF EXISTS uq_storage_usage_samples_storage_id_sampled_on;

ALTER TABLE storage_usage_samples
    DROP CONSTRAINT IF EXISTS fk_storage_usage_samples_storage_id;

DROP TABLE IF EXISTS storage_usage_samples;
-- +goose StatementEnd
