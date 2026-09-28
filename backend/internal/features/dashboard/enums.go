package dashboard

type StorageSpaceStatus string

const (
	StorageSpaceStatusAvailable   StorageSpaceStatus = "AVAILABLE"
	StorageSpaceStatusUnavailable StorageSpaceStatus = "UNAVAILABLE"
	StorageSpaceStatusError       StorageSpaceStatus = "ERROR"
)

type StorageFullForecastStatus string

const (
	StorageFullForecastStatusNotApplicable StorageFullForecastStatus = "NOT_APPLICABLE"
	StorageFullForecastStatusCollecting    StorageFullForecastStatus = "COLLECTING"
	StorageFullForecastStatusFillingUp     StorageFullForecastStatus = "FILLING_UP"
	StorageFullForecastStatusNotFillingUp  StorageFullForecastStatus = "NOT_FILLING_UP"
)
