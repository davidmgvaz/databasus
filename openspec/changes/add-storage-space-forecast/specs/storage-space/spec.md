## Purpose

Defines how Databasus learns the space of each storage, how the dashboard shows it, and how it predicts when a storage will be full.

## ADDED Requirements

### Requirement: Storages report their space when their provider can

The system SHALL read total, used and free bytes from local, NAS, SFTP, rclone and Google Drive storages. A storage whose provider has no way to report space, or whose server refuses the query, SHALL be reported as not reporting space rather than as failed. A storage whose query fails for another reason SHALL be reported as failed with the reason. One storage's failure SHALL NOT prevent the others from being reported.

#### Scenario: A local storage

- **WHEN** the dashboard lists a workspace with a local storage
- **THEN** that storage shows the total, used and free space of the disk that holds the backups

#### Scenario: An S3 storage

- **WHEN** the dashboard lists an S3 storage
- **THEN** that storage shows that its type does not report space

#### Scenario: An unreachable SFTP server

- **WHEN** the SFTP server of a storage cannot be reached
- **THEN** that storage shows that its space could not be read, and the other storages still show theirs

### Requirement: The dashboard lists the workspace's storages above its databases

The system SHALL list every storage of the workspace with its type, the number of the workspace's databases that back up to it, the size of their stored backups, and its used, free and total space with the used share. The list SHALL be collapsible and SHALL remember whether it was collapsed. It SHALL refresh on demand rather than on the dashboard's minute cycle, because reading remote space can take seconds. A reading, including a failed one, MAY be reused for up to one minute.

#### Scenario: A storage nearly full

- **WHEN** a storage has used 90% or more of its space
- **THEN** its usage bar is shown in red

#### Scenario: A collapsed list

- **WHEN** a user collapses the storages list and reloads the page
- **THEN** the list stays collapsed

### Requirement: The space tiles compare the workspace with the installation

The system SHALL show the workspace's total backup size and the free space of its storages. For global admins each tile SHALL also show the value across all workspaces as "this workspace / all workspaces". Free space SHALL count all local storages once, because they share one disk, and every other reporting storage once.

#### Scenario: An admin with two local storages

- **WHEN** an admin's workspace has two local storages and one reporting SFTP storage
- **THEN** the space-left tile adds the local disk's free space once and the SFTP storage's free space once

### Requirement: Only members see a workspace's storages, and only admins see the installation's space

The system SHALL let any member of a workspace, viewers included, and any global admin read the workspace's storages and their space, and SHALL refuse anyone else with the response the databases list gives. The free space across every storage SHALL be returned only to global admins.

#### Scenario: A non-member

- **WHEN** a user outside the workspace requests its storages
- **THEN** the request is refused with "insufficient permissions to access this workspace"

#### Scenario: A member refreshes the storages repeatedly

- **WHEN** a member refreshes the storages list many times within a minute
- **THEN** each storage is contacted at most once in that minute, and no more than four storages are contacted at the same time across all requests

#### Scenario: A member asks for the installation's space

- **WHEN** a user who is not a global admin requests the free space across every storage
- **THEN** the request is refused as forbidden

### Requirement: The system predicts when a storage will be full

The system SHALL record the space of every reporting storage at most once per calendar day in UTC, even when two recordings race, and keep 90 days of records. From the records of the last 30 days it SHALL fit a straight line to the used share over time and show the date the line reaches 100%. With fewer than seven records it SHALL say it is still collecting data. When the used share is flat or falling, or rises so slowly that the full date lies more than ten years away, it SHALL say the storage is not filling up.

#### Scenario: A steadily filling storage

- **WHEN** a storage has eight daily records whose used share rises steadily
- **THEN** the dashboard shows a full date in the future

#### Scenario: A new storage

- **WHEN** a storage has three daily records
- **THEN** the dashboard says it is collecting data, 3 of 7 days

#### Scenario: A barely growing storage

- **WHEN** a storage's used share rises so slowly across its records that it would take about thirty years to fill
- **THEN** the dashboard says it is not filling up

#### Scenario: A storage being cleaned up

- **WHEN** a storage's used share falls across its records
- **THEN** the dashboard says it is not filling up
