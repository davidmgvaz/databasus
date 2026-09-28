## Why

The dashboard shows how much each database stores but not how much room its storage has left. An operator learns that a storage is full when backups start failing. Proxmox Backup Server shows both the free space of each datastore and a date when it will be full; this change brings the same to Databasus for the storages that can report their space.

## What Changes

- Storage providers that can measure their space report total, used and free bytes: local disk, NAS (SMB), SFTP servers with the statvfs extension, rclone remotes that support it, and Google Drive accounts with a quota. S3, Azure Blob and FTP cannot, and say so instead of failing.
- The dashboard gains a collapsible storages table above the databases table: per storage, its type, the databases using it, the size of their backups, used, free and total space, and a usage bar that turns red from 90%.
- The size tile shows "this workspace / all workspaces" for admins, like the count tiles. A new "Space left" tile shows the free space of the workspace's storages, and for admins also across every storage.
- A daily job records one usage sample per reporting storage. After seven daily samples, the dashboard shows when each storage is expected to be full, from a straight-line fit over the last 30 days.
- New routes `GET /api/v1/dashboard/storages?workspace_id=` for workspace members and `GET /api/v1/dashboard/installation/storages` for admins. A new table `storage_usage_samples`.

Nothing existing breaks, so nothing here is BREAKING. Two existing things change:
- The dashboard's fourth tile shows free space instead of the backup size across all workspaces, which moves into the size tile.
- Listing every storage for background work now loads each storage's provider settings, which probing needs.

## Capabilities

### New Capabilities

- `storage-space`: what each storage reports about its space, how the dashboard shows it, and how the full date is predicted.

### Modified Capabilities

None.

## Impact

- Backend: a usage method on five storage providers, parallel probing in the storage service, two dashboard routes, a samples table, a daily sampling job and a line-fit helper.
- Frontend: a storages section on the dashboard, reworked tiles, a TB step in size formatting, new keys in all six dictionaries.
- Answers to `backend/AGENTS.md` and `frontend/AGENTS.md` on top of the root `AGENTS.md`.

## Out of scope

- Showing space anywhere except the dashboard.
- Alerts or notifications when a storage is nearly full.
- Charts of the usage history.
- Quotas entered by hand for storages that cannot report their space.
