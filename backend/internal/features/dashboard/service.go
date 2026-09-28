package dashboard

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"

	backups_core_logical "databasus-backend/internal/features/backups/backups/core/logical"
	physical_core_service "databasus-backend/internal/features/backups/backups/core/physical/service"
	backups_services "databasus-backend/internal/features/backups/backups/services"
	backups_config_logical "databasus-backend/internal/features/backups/config/logical"
	backups_config_physical "databasus-backend/internal/features/backups/config/physical"
	"databasus-backend/internal/features/databases"
	healthcheck_attempt "databasus-backend/internal/features/healthcheck/attempt"
	healthcheck_config "databasus-backend/internal/features/healthcheck/config"
	"databasus-backend/internal/features/storages"
	users_models "databasus-backend/internal/features/users/models"
)

const recentHealthcheckAttemptsLimit = 10

// Errors after the workspace check come from aggregate queries whose text names tables and
// drivers, so the controller answers them with a generic message instead of the text.
var ErrDashboardUnavailable = errors.New("dashboard is unavailable")

type DashboardService struct {
	databaseService             *databases.DatabaseService
	logicalBackupService        *backups_services.LogicalBackupService
	physicalBackupService       *physical_core_service.PhysicalBackupService
	logicalBackupConfigService  *backups_config_logical.BackupConfigService
	physicalBackupConfigService *backups_config_physical.BackupConfigService
	healthcheckConfigService    *healthcheck_config.HealthcheckConfigService
	healthcheckAttemptService   *healthcheck_attempt.HealthcheckAttemptService
	storageService              *storages.StorageService
}

func (s *DashboardService) GetWorkspaceDashboard(
	ctx context.Context,
	user *users_models.User,
	workspaceID uuid.UUID,
) (*WorkspaceDashboard, error) {
	workspaceDatabases, err := s.databaseService.GetDatabasesByWorkspace(ctx, user, workspaceID)
	if err != nil {
		return nil, err
	}

	logicalDatabaseIDs, physicalDatabaseIDs := splitDatabaseIDsByBackupKind(workspaceDatabases)

	logicalTotals, err := s.logicalBackupService.GetBackupTotalsByDatabaseIDs(logicalDatabaseIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	physicalTotals, err := s.physicalBackupService.GetBackupTotalsByDatabaseIDs(physicalDatabaseIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	storagesByDatabaseID, err := s.getStoragesByDatabaseID(logicalDatabaseIDs, physicalDatabaseIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	attemptsByDatabaseID, err := s.getRecentAttemptsOfMonitoredDatabases(
		slices.Concat(logicalDatabaseIDs, physicalDatabaseIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	dashboard := &WorkspaceDashboard{
		Databases: make([]DashboardDatabase, 0, len(workspaceDatabases)),
	}

	for _, database := range workspaceDatabases {
		dashboardDatabase := DashboardDatabase{
			ID:                        database.ID,
			Name:                      database.Name,
			Type:                      database.Type,
			HealthStatus:              database.HealthStatus,
			LastBackupTime:            database.LastBackupTime,
			LastBackupErrorMessage:    database.LastBackupErrorMessage,
			Storage:                   storagesByDatabaseID[database.ID],
			RecentHealthcheckAttempts: attemptsByDatabaseID[database.ID],
		}

		if dashboardDatabase.RecentHealthcheckAttempts == nil {
			dashboardDatabase.RecentHealthcheckAttempts = []DashboardHealthcheckAttempt{}
		}

		if database.Type == databases.DatabaseTypePostgresPhysical {
			applyPhysicalBackupTotals(&dashboardDatabase, physicalTotals[database.ID])
		} else {
			applyLogicalBackupTotals(&dashboardDatabase, logicalTotals[database.ID])
		}

		dashboard.Databases = append(dashboard.Databases, dashboardDatabase)
		dashboard.Totals.DatabasesCount++
		dashboard.Totals.BackupsCount += dashboardDatabase.BackupsCount
		dashboard.Totals.TotalBackupSizeMb += dashboardDatabase.TotalBackupSizeMb
	}

	return dashboard, nil
}

func (s *DashboardService) GetInstallationDashboard() (*DashboardTotals, error) {
	databasesCount, err := s.databaseService.CountDatabases()
	if err != nil {
		return nil, fmt.Errorf("%w: count databases of the installation: %w", ErrDashboardUnavailable, err)
	}

	logicalTotals, err := s.logicalBackupService.GetInstallationBackupTotals()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	physicalTotals, err := s.physicalBackupService.GetInstallationBackupTotals()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	return &DashboardTotals{
		DatabasesCount: databasesCount,
		BackupsCount:   logicalTotals.BackupsCount + physicalTotals.BackupsCount,
		TotalBackupSizeMb: logicalTotals.CompletedBackupSizeMb +
			physicalTotals.CompletedBackupSizeMb +
			physicalTotals.WalSizeMb,
	}, nil
}

func (s *DashboardService) GetWorkspaceStorages(
	ctx context.Context,
	user *users_models.User,
	workspaceID uuid.UUID,
) (*WorkspaceStorages, error) {
	workspaceDashboard, err := s.GetWorkspaceDashboard(ctx, user, workspaceID)
	if err != nil {
		return nil, err
	}

	usageReports, err := s.storageService.GetStorageUsagesByWorkspace(ctx, user, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	databasesCountByStorageID := make(map[uuid.UUID]int64)
	backupsSizeMbByStorageID := make(map[uuid.UUID]float64)
	for _, dashboardDatabase := range workspaceDashboard.Databases {
		if dashboardDatabase.Storage == nil {
			continue
		}

		databasesCountByStorageID[dashboardDatabase.Storage.ID]++
		backupsSizeMbByStorageID[dashboardDatabase.Storage.ID] += dashboardDatabase.TotalBackupSizeMb
	}

	workspaceStorages := &WorkspaceStorages{
		Storages:       make([]DashboardStorageUsage, 0, len(usageReports)),
		FreeSpaceBytes: sumFreeSpaceBytes(usageReports),
	}

	for _, usageReport := range usageReports {
		storageUsage := toDashboardStorageUsage(usageReport)
		storageUsage.DatabasesCount = databasesCountByStorageID[usageReport.StorageID]
		storageUsage.BackupsSizeMb = backupsSizeMbByStorageID[usageReport.StorageID]

		workspaceStorages.Storages = append(workspaceStorages.Storages, storageUsage)
	}

	return workspaceStorages, nil
}

func (s *DashboardService) GetInstallationStorages(ctx context.Context) (*InstallationStorages, error) {
	usageReports, err := s.storageService.GetAllStorageUsages(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDashboardUnavailable, err)
	}

	return &InstallationStorages{FreeSpaceBytes: sumFreeSpaceBytes(usageReports)}, nil
}

func (s *DashboardService) getStoragesByDatabaseID(
	logicalDatabaseIDs []uuid.UUID,
	physicalDatabaseIDs []uuid.UUID,
) (map[uuid.UUID]*DashboardStorage, error) {
	storagesByDatabaseID := make(map[uuid.UUID]*DashboardStorage)

	logicalBackupConfigs, err := s.logicalBackupConfigService.GetBackupConfigsByDatabaseIDs(logicalDatabaseIDs)
	if err != nil {
		return nil, err
	}

	for _, backupConfig := range logicalBackupConfigs {
		storagesByDatabaseID[backupConfig.DatabaseID] = toDashboardStorage(backupConfig.Storage)
	}

	physicalBackupConfigs, err := s.physicalBackupConfigService.GetBackupConfigsByDatabaseIDs(physicalDatabaseIDs)
	if err != nil {
		return nil, err
	}

	for _, backupConfig := range physicalBackupConfigs {
		storagesByDatabaseID[backupConfig.DatabaseID] = toDashboardStorage(backupConfig.Storage)
	}

	return storagesByDatabaseID, nil
}

func (s *DashboardService) getRecentAttemptsOfMonitoredDatabases(
	databaseIDs []uuid.UUID,
) (map[uuid.UUID][]DashboardHealthcheckAttempt, error) {
	healthcheckConfigs, err := s.healthcheckConfigService.GetConfigsByDatabaseIDs(databaseIDs)
	if err != nil {
		return nil, err
	}

	monitoredDatabaseIDs := make([]uuid.UUID, 0, len(healthcheckConfigs))
	for _, healthcheckConfig := range healthcheckConfigs {
		if healthcheckConfig.IsHealthcheckEnabled {
			monitoredDatabaseIDs = append(monitoredDatabaseIDs, healthcheckConfig.DatabaseID)
		}
	}

	attemptsByDatabaseID, err := s.healthcheckAttemptService.GetRecentAttemptsByDatabaseIDs(
		monitoredDatabaseIDs,
		recentHealthcheckAttemptsLimit,
	)
	if err != nil {
		return nil, err
	}

	dashboardAttemptsByDatabaseID := make(map[uuid.UUID][]DashboardHealthcheckAttempt, len(attemptsByDatabaseID))
	for databaseID, attempts := range attemptsByDatabaseID {
		dashboardAttempts := make([]DashboardHealthcheckAttempt, 0, len(attempts))
		for _, attempt := range attempts {
			dashboardAttempts = append(dashboardAttempts, DashboardHealthcheckAttempt{
				Status:    attempt.Status,
				CreatedAt: attempt.CreatedAt,
			})
		}

		dashboardAttemptsByDatabaseID[databaseID] = dashboardAttempts
	}

	return dashboardAttemptsByDatabaseID, nil
}

func splitDatabaseIDsByBackupKind(
	workspaceDatabases []*databases.Database,
) (logicalDatabaseIDs, physicalDatabaseIDs []uuid.UUID) {
	logicalDatabaseIDs = []uuid.UUID{}
	physicalDatabaseIDs = []uuid.UUID{}

	for _, database := range workspaceDatabases {
		if database.Type == databases.DatabaseTypePostgresPhysical {
			physicalDatabaseIDs = append(physicalDatabaseIDs, database.ID)
		} else {
			logicalDatabaseIDs = append(logicalDatabaseIDs, database.ID)
		}
	}

	return logicalDatabaseIDs, physicalDatabaseIDs
}

func applyLogicalBackupTotals(
	dashboardDatabase *DashboardDatabase,
	totals backups_core_logical.DatabaseBackupTotals,
) {
	dashboardDatabase.BackupsCount = totals.BackupsCount
	dashboardDatabase.CompletedBackupsCount = totals.CompletedBackupsCount
	dashboardDatabase.FailedBackupsCount = totals.FailedBackupsCount
	dashboardDatabase.MeanBackupSizeMb = getMeanSizeMb(totals.CompletedBackupSizeMb, totals.CompletedBackupsCount)
	dashboardDatabase.TotalBackupSizeMb = totals.CompletedBackupSizeMb
}

func applyPhysicalBackupTotals(
	dashboardDatabase *DashboardDatabase,
	totals physical_core_service.DatabasePhysicalBackupTotals,
) {
	dashboardDatabase.BackupsCount = totals.BackupsCount
	dashboardDatabase.CompletedBackupsCount = totals.CompletedBackupsCount
	dashboardDatabase.FailedBackupsCount = totals.FailedBackupsCount
	dashboardDatabase.MeanBackupSizeMb = getMeanSizeMb(totals.CompletedBackupSizeMb, totals.CompletedBackupsCount)
	dashboardDatabase.TotalBackupSizeMb = totals.CompletedBackupSizeMb + totals.WalSizeMb
}

func getMeanSizeMb(totalSizeMb float64, backupsCount int64) *float64 {
	if backupsCount == 0 {
		return nil
	}

	return new(totalSizeMb / float64(backupsCount))
}

func toDashboardStorage(storage *storages.Storage) *DashboardStorage {
	if storage == nil {
		return nil
	}

	return &DashboardStorage{
		ID:   storage.ID,
		Name: storage.Name,
		Type: storage.Type,
	}
}

func toDashboardStorageUsage(usageReport storages.StorageUsageReport) DashboardStorageUsage {
	storageUsage := DashboardStorageUsage{
		ID:   usageReport.StorageID,
		Name: usageReport.StorageName,
		Type: usageReport.StorageType,
	}

	switch {
	case usageReport.Usage != nil:
		storageUsage.SpaceStatus = StorageSpaceStatusAvailable
		storageUsage.Space = &DashboardStorageSpace{
			TotalBytes: usageReport.Usage.TotalBytes,
			UsedBytes:  usageReport.Usage.UsedBytes,
			FreeBytes:  usageReport.Usage.FreeBytes,
		}
	case usageReport.IsUnavailable:
		storageUsage.SpaceStatus = StorageSpaceStatusUnavailable
	default:
		storageUsage.SpaceStatus = StorageSpaceStatusError
		if usageReport.ProbeError != nil {
			storageUsage.SpaceErrorMessage = new(usageReport.ProbeError.Error())
		}
	}

	return storageUsage
}

// Every LOCAL storage writes to the same disk, so its free space counts once. Remote storages
// count once each, because the dashboard cannot tell two of them sharing a volume apart.
func sumFreeSpaceBytes(usageReports []storages.StorageUsageReport) *int64 {
	var freeSpaceBytes *int64
	isLocalDiskCounted := false

	for _, usageReport := range usageReports {
		if usageReport.Usage == nil {
			continue
		}

		if usageReport.StorageType == storages.StorageTypeLocal {
			if isLocalDiskCounted {
				continue
			}

			isLocalDiskCounted = true
		}

		if freeSpaceBytes == nil {
			freeSpaceBytes = new(int64(0))
		}

		*freeSpaceBytes += usageReport.Usage.FreeBytes
	}

	return freeSpaceBytes
}
