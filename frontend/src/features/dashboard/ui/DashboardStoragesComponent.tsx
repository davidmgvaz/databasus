import {
  DownOutlined,
  ExclamationCircleOutlined,
  ReloadOutlined,
  RightOutlined,
} from '@ant-design/icons';
import { Button, Progress, Spin, Table, Tooltip } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

import { STORAGE_TYPE_LABEL_KEYS, getStorageLogoFromType } from '../../../entity/storages';
import { usePersistentState } from '../../../shared/hooks';
import { useLocale } from '../../../shared/i18n';
import { formatSizeBytes, formatSizeMb } from '../../../shared/lib';
import type { DashboardStorageUsage } from '../model/DashboardStorageUsage';
import { StorageSpaceStatus } from '../model/StorageSpaceStatus';

interface Props {
  storageUsages?: DashboardStorageUsage[];
  isLoading: boolean;
  onRefresh: () => void;
}

// eslint-disable-next-line i18next/no-literal-string -- localStorage key
const EXPANDED_STORAGE_KEY = 'dashboard_storages_expanded';
const NEARLY_FULL_PERCENT = 90;
const NO_VALUE = '-';

const getUsedPercent = (storageUsage: DashboardStorageUsage): number | undefined => {
  if (!storageUsage.space || storageUsage.space.totalBytes <= 0) {
    return undefined;
  }

  return (storageUsage.space.usedBytes / storageUsage.space.totalBytes) * 100;
};

const renderNoValue = () => <span className="text-gray-400 dark:text-gray-500">{NO_VALUE}</span>;

const renderStorageName = (storageUsage: DashboardStorageUsage) => (
  <span className="inline-flex min-w-0 items-center gap-2">
    <img src={getStorageLogoFromType(storageUsage.type)} alt="" className="h-4 w-4 shrink-0" />
    <span className="font-medium break-all">{storageUsage.name}</span>
  </span>
);

const renderMobileField = (label: string, value: ReactNode) => (
  <div>
    <div className="text-xs text-gray-500 dark:text-gray-400">{label}</div>
    <div className="text-sm">{value}</div>
  </div>
);

export const DashboardStoragesComponent = ({ storageUsages, isLoading, onRefresh }: Props) => {
  const { t } = useTranslation();
  const { formatNumber } = useLocale();

  const [isExpanded, setIsExpanded] = usePersistentState(EXPANDED_STORAGE_KEY, true);

  const renderSpaceBytes = (storageUsage: DashboardStorageUsage, spaceBytes?: number) =>
    storageUsage.space && spaceBytes !== undefined
      ? formatSizeBytes(spaceBytes, formatNumber)
      : renderNoValue();

  const renderUsage = (storageUsage: DashboardStorageUsage) => {
    if (storageUsage.spaceStatus === StorageSpaceStatus.UNAVAILABLE) {
      return (
        <span className="text-xs text-gray-400 dark:text-gray-500">
          {t('dashboard.storages.spaceUnavailable')}
        </span>
      );
    }

    if (storageUsage.spaceStatus === StorageSpaceStatus.ERROR) {
      return (
        <Tooltip title={storageUsage.spaceErrorMessage}>
          <span className="inline-flex items-center gap-1 text-xs text-red-600 dark:text-red-400">
            <ExclamationCircleOutlined />
            {t('dashboard.storages.spaceError')}
          </span>
        </Tooltip>
      );
    }

    const usedPercent = getUsedPercent(storageUsage);
    if (usedPercent === undefined) {
      return renderNoValue();
    }

    const isNearlyFull = usedPercent >= NEARLY_FULL_PERCENT;

    return (
      <div className="min-w-[120px]">
        <Progress
          percent={Number(usedPercent.toFixed(1))}
          size="small"
          strokeColor={isNearlyFull ? '#ef4444' : undefined}
          format={(percent) => `${formatNumber(percent ?? 0)}%`}
        />
      </div>
    );
  };

  const columns: ColumnsType<DashboardStorageUsage> = [
    {
      title: t('dashboard.storages.columns.storage'),
      key: 'name',
      render: (_: unknown, storageUsage: DashboardStorageUsage) => renderStorageName(storageUsage),
      sorter: (a, b) => a.name.localeCompare(b.name),
    },
    {
      title: t('dashboard.storages.columns.type'),
      key: 'type',
      render: (_: unknown, storageUsage: DashboardStorageUsage) =>
        t(STORAGE_TYPE_LABEL_KEYS[storageUsage.type]),
    },
    {
      title: t('dashboard.storages.columns.databases'),
      key: 'databasesCount',
      render: (_: unknown, storageUsage: DashboardStorageUsage) =>
        formatNumber(storageUsage.databasesCount),
      sorter: (a, b) => a.databasesCount - b.databasesCount,
    },
    {
      title: t('dashboard.storages.columns.backupsSize'),
      key: 'backupsSizeMb',
      render: (_: unknown, storageUsage: DashboardStorageUsage) =>
        formatSizeMb(storageUsage.backupsSizeMb, formatNumber),
      sorter: (a, b) => a.backupsSizeMb - b.backupsSizeMb,
    },
    {
      title: t('dashboard.storages.columns.used'),
      key: 'usedBytes',
      render: (_: unknown, storageUsage: DashboardStorageUsage) =>
        renderSpaceBytes(storageUsage, storageUsage.space?.usedBytes),
    },
    {
      title: t('dashboard.storages.columns.free'),
      key: 'freeBytes',
      render: (_: unknown, storageUsage: DashboardStorageUsage) =>
        renderSpaceBytes(storageUsage, storageUsage.space?.freeBytes),
      sorter: (a, b) => (a.space?.freeBytes ?? 0) - (b.space?.freeBytes ?? 0),
    },
    {
      title: t('dashboard.storages.columns.total'),
      key: 'totalBytes',
      render: (_: unknown, storageUsage: DashboardStorageUsage) =>
        renderSpaceBytes(storageUsage, storageUsage.space?.totalBytes),
    },
    {
      title: t('dashboard.storages.columns.usage'),
      key: 'usage',
      render: (_: unknown, storageUsage: DashboardStorageUsage) => renderUsage(storageUsage),
      sorter: (a, b) => (getUsedPercent(a) ?? -1) - (getUsedPercent(b) ?? -1),
    },
  ];

  const renderContent = () => {
    if (isLoading && !storageUsages) {
      return (
        <div className="flex justify-center py-6">
          <Spin />
        </div>
      );
    }

    if (!storageUsages || storageUsages.length === 0) {
      return (
        <div className="mt-3 text-sm text-gray-500 dark:text-gray-400">
          {t('dashboard.storages.empty')}
        </div>
      );
    }

    return (
      <>
        <div className="mt-4 hidden md:block">
          <Table
            bordered
            columns={columns}
            dataSource={storageUsages}
            rowKey="id"
            size="small"
            pagination={false}
            loading={isLoading}
          />
        </div>

        <div className="mt-3 md:hidden">
          {storageUsages.map((storageUsage) => (
            <div
              key={storageUsage.id}
              className="mb-2 rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800"
            >
              <div className="flex items-start justify-between gap-2">
                {renderStorageName(storageUsage)}
                <span className="shrink-0 text-xs text-gray-500 dark:text-gray-400">
                  {t(STORAGE_TYPE_LABEL_KEYS[storageUsage.type])}
                </span>
              </div>

              <div className="mt-2">{renderUsage(storageUsage)}</div>

              <div className="mt-3 grid grid-cols-2 gap-3">
                {renderMobileField(
                  t('dashboard.storages.columns.databases'),
                  formatNumber(storageUsage.databasesCount),
                )}
                {renderMobileField(
                  t('dashboard.storages.columns.backupsSize'),
                  formatSizeMb(storageUsage.backupsSizeMb, formatNumber),
                )}
                {renderMobileField(
                  t('dashboard.storages.columns.free'),
                  renderSpaceBytes(storageUsage, storageUsage.space?.freeBytes),
                )}
                {renderMobileField(
                  t('dashboard.storages.columns.total'),
                  renderSpaceBytes(storageUsage, storageUsage.space?.totalBytes),
                )}
              </div>
            </div>
          ))}
        </div>
      </>
    );
  };

  return (
    <div className="mt-2 rounded bg-white p-3 shadow md:mt-3 md:p-5 dark:bg-gray-800">
      <div className="flex items-center justify-between gap-2">
        <button
          type="button"
          className="flex cursor-pointer items-center gap-2 text-left"
          onClick={() => setIsExpanded(!isExpanded)}
          aria-expanded={isExpanded}
        >
          {isExpanded ? <DownOutlined /> : <RightOutlined />}
          <h2 className="text-lg font-bold md:text-xl">{t('dashboard.storages.title')}</h2>
        </button>

        <Tooltip title={t('common.actions.refresh')}>
          <Button
            type="text"
            icon={<ReloadOutlined spin={isLoading} />}
            onClick={onRefresh}
            disabled={isLoading}
            aria-label={t('common.actions.refresh')}
          />
        </Tooltip>
      </div>

      {isExpanded && renderContent()}
    </div>
  );
};
