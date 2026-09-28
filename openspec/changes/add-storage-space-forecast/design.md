## Context

The dashboard (`openspec/changes/add-workspace-dashboard/`) aggregates backups per database. Storages are opened through one provider struct per type behind `StorageFileSaver` (`backend/internal/features/storages/interfaces.go:13`), dispatched by `(*Storage).getSpecificStorage` (`backend/internal/features/storages/model.go`). No provider reports its space; the only disk measurement is the host disk behind `GET /disk/usage` (`backend/internal/features/disk/service.go:15`). The libraries the providers already use expose space queries: gopsutil `disk.Usage`, go-smb2 `(*Share).Statfs`, pkg/sftp `(*Client).StatVFS`, rclone `Features().About`, and Drive `About.Get` with `storageQuota`.

## Goals / Non-Goals

**Goals:**
- Free space per storage on the dashboard, for every provider that can report it.
- A full-date forecast from daily records.

**Non-Goals:**
- Space outside the dashboard, alerts, history charts and hand-entered quotas (see the proposal).

## Decisions

### An optional interface instead of a new required method

`StorageUsageReporter` is a separate interface that `(*Storage).GetUsage` type-asserts. Providers without a space API do not implement it and are reported as "not reporting space".

Rejected alternative: adding `GetUsage` to `StorageFileSaver`. S3, Azure Blob and FTP would each carry a stub that only returns "unavailable".

### A leaf package for the usage type

`storages/space` holds `Usage` and `ErrUsageUnavailable`. The provider packages under `storages/models` cannot import `storages`, which imports them.

Rejected alternative: returning three bare integers from each provider. Same-typed adjacent results invite transposed use, which the root `AGENTS.md` sends to a named struct.

### Reuse each provider's own connection setup

Each `GetUsage` opens its client the way that provider's `TestConnection` does (`createSessionWithContext` for NAS, `connectWithContext` for SFTP, `getFs` for rclone, `withRetryOnAuth` for Drive), so timeouts, decryption and authentication stay in one place per provider.

### Measuring the local disk at the nearest existing folder

The backups folder may not exist before the first backup, so the local provider walks up to the nearest existing parent, which lives on the same disk.

Rejected alternative: measuring the parent folder unconditionally, as `GET /disk/usage` does. It fails the same way on a fresh install whose data folder has not been created yet.

### Probe in parallel, never fail the listing

`StorageService.probeStorageUsages` probes every storage concurrently under a 15-second timeout each and records a failure per storage. One hung remote therefore costs at most 15 seconds and never hides the space of the others.

Rejected alternative: failing the request on the first probe error. A single offline NAS would blank the whole storages list.

### Separate routes from the minute-refreshed dashboard

`GET /dashboard/storages` and `GET /dashboard/installation/storages` are separate from `GET /dashboard` and `GET /dashboard/installation`. The page loads them on mount and on a refresh button.

Rejected alternative: adding space to the existing routes. Their one-minute refresh would probe every remote storage every minute.

### Counting free space

All local storages write to one disk, so the free-space sums count it once. Remote storages count once each; two remotes on one volume cannot be told apart, and the tile hint says what is counted.

### Daily records and a line fit for the forecast

A background job records one sample per reporting storage per day into `storage_usage_samples`, in the shape of the telemetry job (`backend/internal/features/telemetry/background_service.go`). The forecast fits a least-squares line to used share against time over 30 days, needs seven samples, and projects to 100%, as Proxmox Backup Server does.

Rejected alternative: computing the rate from backup sizes in the catalog. Storages also hold files Databasus does not know about, so only a reading of the storage itself tracks when it fills.

## AGENTS.md constraints

- `backend/AGENTS.md`: services, not repositories, cross feature lines; migrations declare tables, constraints and indexes separately; background jobs panic on a second `Run` and log `job_id` and `job_name`.
- `frontend/AGENTS.md`: the storages section stays in `features/dashboard`, its only consumer; every string is a dictionary key in all six languages.

## Risks / Trade-offs

- **Google Drive scope.** The provider requests the `drive.file` scope. If Google refuses `about.get` under it, the answer is 403, which the provider maps to "not reporting space".
- **Remote latency.** Probing a slow remote delays the storages list up to 15 seconds; the rest of the dashboard is unaffected.
- **Forecast accuracy.** A straight line ignores retention cycles. It gives the same rough guidance as Proxmox, not a guarantee.
