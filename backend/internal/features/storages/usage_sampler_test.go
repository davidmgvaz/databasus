package storages

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	users_enums "databasus-backend/internal/features/users/enums"
	users_testing "databasus-backend/internal/features/users/testing"
	workspaces_controllers "databasus-backend/internal/features/workspaces/controllers"
	workspaces_testing "databasus-backend/internal/features/workspaces/testing"
	"databasus-backend/internal/util/logger"
)

func Test_RecordUsageSamples_RunTwiceInOneDay_StoresOneSamplePerLocalStorage(t *testing.T) {
	router := workspaces_testing.CreateTestRouter(
		workspaces_controllers.GetWorkspaceController(),
		workspaces_controllers.GetMembershipController(),
	)
	owner := users_testing.CreateTestUser(t.Context(), users_enums.UserRoleMember)
	workspace := workspaces_testing.CreateTestWorkspace(t.Context(), "Usage Sampling Workspace", owner, router)
	storage := CreateTestStorage(workspace.ID)
	t.Cleanup(func() {
		RemoveTestStorage(t.Context(), storage.ID)
		workspaces_testing.RemoveTestWorkspace(t.Context(), workspace, router)
	})

	sampler := &StorageUsageSampler{
		storageService:        storageService,
		usageSampleRepository: storageUsageSampleRepository,
		logger:                logger.GetLogger(),
	}

	sampler.recordUsageSamples(t.Context())
	sampler.recordUsageSamples(t.Context())

	samplesByStorageID, err := storageService.GetUsageSamplesSince(
		[]uuid.UUID{storage.ID},
		time.Now().UTC().Add(-time.Hour),
	)
	require.NoError(t, err)

	storageSamples := samplesByStorageID[storage.ID]
	require.Len(t, storageSamples, 1)
	assert.Positive(t, storageSamples[0].TotalBytes)
	assert.LessOrEqual(t, storageSamples[0].FreeBytes, storageSamples[0].TotalBytes)
	assert.LessOrEqual(t, storageSamples[0].UsedBytes, storageSamples[0].TotalBytes)
}
