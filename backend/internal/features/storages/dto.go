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
// Errors do not survive JSON, so a cached failure keeps only its message.
type cachedStorageUsage struct {
	Usage             *storage_space.Usage `json:"usage,omitempty"`
	IsUnavailable     bool                 `json:"isUnavailable"`
	ProbeErrorMessage string               `json:"probeErrorMessage,omitempty"`
}

type StorageUsageReport struct {
	StorageID     uuid.UUID
	StorageName   string
	StorageType   StorageType
	Usage         *storage_space.Usage
	IsUnavailable bool
	ProbeError    error
}
