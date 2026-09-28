package dashboard

type StorageSpaceStatus string

const (
	StorageSpaceStatusAvailable   StorageSpaceStatus = "AVAILABLE"
	StorageSpaceStatusUnavailable StorageSpaceStatus = "UNAVAILABLE"
	StorageSpaceStatusError       StorageSpaceStatus = "ERROR"
)
