import { getApplicationServer } from '../../../constants';
import RequestOptions from '../../../shared/api/RequestOptions';
import { apiHelper } from '../../../shared/api/apiHelper';
import type { DashboardTotals } from '../model/DashboardTotals';
import type { InstallationStorages } from '../model/InstallationStorages';
import type { WorkspaceDashboard } from '../model/WorkspaceDashboard';
import type { WorkspaceStorages } from '../model/WorkspaceStorages';

export const dashboardApi = {
  async getWorkspaceDashboard(workspaceId: string) {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchGetJson<WorkspaceDashboard>(
      `${getApplicationServer()}/api/v1/dashboard?workspace_id=${workspaceId}`,
      requestOptions,
      true,
    );
  },

  async getWorkspaceStorages(workspaceId: string) {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchGetJson<WorkspaceStorages>(
      `${getApplicationServer()}/api/v1/dashboard/storages?workspace_id=${workspaceId}`,
      requestOptions,
      true,
    );
  },

  async getInstallationStorages() {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchGetJson<InstallationStorages>(
      `${getApplicationServer()}/api/v1/dashboard/installation/storages`,
      requestOptions,
      true,
    );
  },

  async getInstallationDashboard() {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchGetJson<DashboardTotals>(
      `${getApplicationServer()}/api/v1/dashboard/installation`,
      requestOptions,
      true,
    );
  },
};
