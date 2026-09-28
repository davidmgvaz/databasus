package storage_space

import "errors"

// Providers without a space API, and servers that refuse the query, report this instead of an
// error, so the dashboard can tell "cannot know" apart from "failed to ask".
var ErrUsageUnavailable = errors.New("storage does not report its space")

type Usage struct {
	TotalBytes int64
	UsedBytes  int64
	FreeBytes  int64
}
