package storages

import (
	"github.com/google/uuid"

	storage_space "databasus-backend/internal/features/storages/space"
)

type TransferStorageRequest struct {
	TargetWorkspaceID uuid.UUID `json:"targetWorkspaceId" binding:"required"`
}

// A failed or unsupported probe is kept per storage instead of failing the whole listing, since
// one unreachable remote must not hide the space of the others.
type StorageUsageReport struct {
	StorageID     uuid.UUID
	StorageName   string
	StorageType   StorageType
	Usage         *storage_space.Usage
	IsUnavailable bool
	ProbeError    error
}
