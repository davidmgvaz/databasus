## 1. Free space

- [x] 1.1 Add the `storages/space` package, the `StorageUsageReporter` interface and `(*Storage).GetUsage`; verify `go build ./...`
- [x] 1.2 Implement `GetUsage` on the local, NAS, SFTP, rclone and Google Drive providers; verify `go vet ./internal/features/storages/...`
- [x] 1.3 Probe storages in parallel in the storage service and load provider settings in `GetAllStorages`; verify `go vet ./internal/features/storages/...`
- [x] 1.4 Add `GET /dashboard/storages` and `GET /dashboard/installation/storages`; verify `make swagger` lists both
- [x] 1.5 Add the dashboard controller tests and the provider `GetUsage` subtest; verify `go vet ./internal/features/dashboard/... ./internal/features/storages/...`
- [ ] 1.6 Run the new backend tests; verify `go test ./internal/features/dashboard/... ./internal/features/storages/... -count=1`
- [x] 1.7 Add the TB step and `formatSizeBytes`, the storages section and the reworked tiles; verify `pnpm test` and `pnpm build`
- [x] 1.8 Add the dictionary keys in all six languages; verify `pnpm build`
- [x] 1.9 Lint; verify `make lint`, `pnpm lint`

## 2. Estimated full

- [ ] 2.1 Add the `storage_usage_samples` migration and repository; verify `make migration-up`
- [x] 2.2 Add the daily sampling job and start it from `cmd/main.go`; verify `go build ./...`
- [x] 2.3 Add the line-fit helper with a unit test; verify `go test ./internal/util/statistics/... -count=1`
- [ ] 2.4 Add the forecast to the storages response with controller tests; verify `go test ./internal/features/dashboard/... -count=1`
- [x] 2.5 Show the forecast column on the dashboard in all six languages; verify `pnpm build`
- [x] 2.6 Lint; verify `make lint`, `pnpm lint`

## 3. Review feedback

- [x] 3.1 Cache each storage's reading for a minute and cap parallel probes at four; verify `make lint`
- [x] 3.2 Enforce one usage sample per storage per UTC day with a unique constraint; verify `go vet ./internal/features/storages/...`
- [x] 3.3 Load the workspace dashboard and installation totals independently; verify `pnpm lint` and `pnpm build`

## 4. Review

- [ ] 4.1 Run the post-implementation compliance review from the root `AGENTS.md` and resolve its findings
