import type { StorageType } from '../../../entity/storages';
import type { DashboardStorageSpace } from './DashboardStorageSpace';
import type { StorageFullForecast } from './StorageFullForecast';
import type { StorageSpaceStatus } from './StorageSpaceStatus';

export interface DashboardStorageUsage {
  id: string;
  name: string;
  type: StorageType;
  databasesCount: number;
  backupsSizeMb: number;
  spaceStatus: StorageSpaceStatus;
  space?: DashboardStorageSpace;
  spaceErrorMessage?: string;
  fullForecast: StorageFullForecast;
}
