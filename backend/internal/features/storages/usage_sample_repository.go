package storages

import (
	"time"

	"github.com/google/uuid"

	db "databasus-backend/internal/storage"
)

type StorageUsageSampleRepository struct{}

func (r *StorageUsageSampleRepository) Insert(sample *StorageUsageSample) error {
	return db.GetDb().Create(sample).Error
}

func (r *StorageUsageSampleRepository) HasSampleSince(storageID uuid.UUID, since time.Time) (bool, error) {
	var count int64

	if err := db.
		GetDb().
		Model(&StorageUsageSample{}).
		Where("storage_id = ? AND sampled_at >= ?", storageID, since).
		Limit(1).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
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
