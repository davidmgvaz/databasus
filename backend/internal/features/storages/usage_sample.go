package storages

import (
	"time"

	"github.com/google/uuid"
)

type StorageUsageSample struct {
	ID         uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	StorageID  uuid.UUID `gorm:"column:storage_id;type:uuid;not null"`
	SampledAt  time.Time `gorm:"column:sampled_at;type:timestamptz;not null"`
	SampledOn  time.Time `gorm:"column:sampled_on;type:date;not null"`
	TotalBytes int64     `gorm:"column:total_bytes;not null"`
	UsedBytes  int64     `gorm:"column:used_bytes;not null"`
	FreeBytes  int64     `gorm:"column:free_bytes;not null"`
}

func (s *StorageUsageSample) TableName() string {
	return "storage_usage_samples"
}
