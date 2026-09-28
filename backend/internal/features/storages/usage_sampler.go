package storages

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	storageUsageSamplerJobName = "storage_usage_sampling"

	storageUsageSamplingWarmup   = 1 * time.Minute
	storageUsageSamplingInterval = 24 * time.Hour

	// A restart inside a day must not add a second point for that day, or the forecast's
	// "days of data" count and its slope would both drift.
	storageUsageSampleMinimumSpacing = 20 * time.Hour
	storageUsageSampleRetention      = 90 * 24 * time.Hour
)

type StorageUsageSampler struct {
	storageService        *StorageService
	usageSampleRepository *StorageUsageSampleRepository
	logger                *slog.Logger
	hasRun                atomic.Bool
}

func (s *StorageUsageSampler) Run(ctx context.Context) {
	if s.hasRun.Swap(true) {
		panic(fmt.Sprintf("%T.Run() called multiple times", s))
	}

	warmupTimer := time.NewTimer(storageUsageSamplingWarmup)
	defer warmupTimer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-warmupTimer.C:
	}

	s.recordUsageSamples(ctx)

	ticker := time.NewTicker(storageUsageSamplingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.recordUsageSamples(ctx)
		}
	}
}

func (s *StorageUsageSampler) recordUsageSamples(ctx context.Context) {
	logger := s.logger.With("job_id", uuid.New(), "job_name", storageUsageSamplerJobName)

	usageReports, err := s.storageService.GetAllStorageUsages(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to read storage usages", "error", err)
		return
	}

	now := time.Now().UTC()
	recordedSamplesCount := 0

	for _, usageReport := range usageReports {
		if usageReport.Usage == nil || usageReport.Usage.TotalBytes <= 0 {
			continue
		}

		hasRecentSample, err := s.usageSampleRepository.HasSampleSince(
			usageReport.StorageID,
			now.Add(-storageUsageSampleMinimumSpacing),
		)
		if err != nil {
			logger.ErrorContext(ctx, "failed to check last storage usage sample",
				"storage_id", usageReport.StorageID, "error", err)
			continue
		}
		if hasRecentSample {
			continue
		}

		if err := s.usageSampleRepository.Insert(&StorageUsageSample{
			StorageID:  usageReport.StorageID,
			SampledAt:  now,
			TotalBytes: usageReport.Usage.TotalBytes,
			UsedBytes:  usageReport.Usage.UsedBytes,
			FreeBytes:  usageReport.Usage.FreeBytes,
		}); err != nil {
			logger.ErrorContext(ctx, "failed to record storage usage sample",
				"storage_id", usageReport.StorageID, "error", err)
			continue
		}

		recordedSamplesCount++
	}

	if err := s.usageSampleRepository.DeleteOlderThan(now.Add(-storageUsageSampleRetention)); err != nil {
		logger.ErrorContext(ctx, "failed to delete old storage usage samples", "error", err)
	}

	logger.InfoContext(ctx, "recorded storage usage samples", "recorded_samples_count", recordedSamplesCount)
}
