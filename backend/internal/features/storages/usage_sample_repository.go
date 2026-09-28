package storages

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	db "databasus-backend/internal/storage"
)

type StorageUsageSampleRepository struct{}

// The unique (storage_id, sampled_on) constraint makes a second sample of the same UTC day a
// no-op even when two processes race, so the forecast never counts a day twice.
func (r *StorageUsageSampleRepository) InsertOncePerDay(sample *StorageUsageSample) (bool, error) {
	sampledAt := sample.SampledAt.UTC()
	sample.SampledOn = time.Date(sampledAt.Year(), sampledAt.Month(), sampledAt.Day(), 0, 0, 0, 0, time.UTC)

	insertResult := db.GetDb().Clauses(clause.OnConflict{DoNothing: true}).Create(sample)
	if insertResult.Error != nil {
		return false, insertResult.Error
	}

	return insertResult.RowsAffected > 0, nil
}

func (r *StorageUsageSampleRepository) FindSinceByStorageIDs(
	storageIDs []uuid.UUID,
	since time.Time,
) ([]StorageUsageSample, error) {
	samples := []StorageUsageSample{}

	if len(storageIDs) == 0 {
		return samples, nil
	}

	if err := db.
		GetDb().
		Where("storage_id IN ? AND sampled_at >= ?", storageIDs, since).
		Order("storage_id, sampled_at").
		Find(&samples).Error; err != nil {
		return nil, err
	}

	return samples, nil
}

func (r *StorageUsageSampleRepository) DeleteOlderThan(before time.Time) error {
	return db.GetDb().Where("sampled_at < ?", before).Delete(&StorageUsageSample{}).Error
}
