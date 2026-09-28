import type { StorageFullForecastStatus } from './StorageFullForecastStatus';

export interface StorageFullForecast {
  status: StorageFullForecastStatus;
  estimatedFullAt?: Date;
  sampleCount: number;
  requiredSampleCount: number;
}
