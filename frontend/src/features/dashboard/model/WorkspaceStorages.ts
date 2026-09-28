import type { DashboardStorageUsage } from './DashboardStorageUsage';

export interface WorkspaceStorages {
  storages: DashboardStorageUsage[];
  freeSpaceBytes?: number;
}
