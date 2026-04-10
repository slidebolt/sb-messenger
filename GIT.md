# Git Workflow for sb-messenger

This repository contains the Slidebolt Messenger service, which wraps the NATS server and provides Slidebolt-specific messaging features. It produces a standalone binary.

## Dependencies
- **Internal:**
  - `sb-contract`: Core interfaces and shared structures.
  - `sb-messenger-sdk`: Shared messaging interfaces and NATS implementation.
  - `sb-runtime`: Core execution environment and logging.
- **External:** 
  - `github.com/nats-io/nats-server/v2`: Embedded NATS server.

## Build Process
- **Type:** Go Application (Service).
- **Consumption:** Run as the primary communication hub for Slidebolt.
- **Artifacts:** Produces a binary named `sb-messenger`.
- **Command:** `go build -o sb-messenger ./cmd/sb-messenger`
- **Validation:** 
  - Validated through unit tests: `go test -v ./...`
  - Validated by successful compilation of the binary.

## Pre-requisites & Publishing
As the messaging hub, `sb-messenger` must be updated whenever the core contracts, runtime, or messaging SDK is changed.

**Before publishing:**
1. Determine current tag: `git tag | sort -V | tail -n 1`
2. Ensure all local tests pass: `go test -v ./...`
3. Ensure the binary builds: `go build -o sb-messenger ./cmd/sb-messenger`

**Publishing Order:**
1. Ensure `sb-contract`, `sb-messenger-sdk`, and `sb-runtime` are tagged and pushed (e.g., `v1.0.4`).
2. Update `sb-messenger/go.mod` to reference the latest tags.
3. Determine next semantic version for `sb-messenger` (e.g., `v1.0.4`).
4. Commit and push the changes to `main`.
5. Tag the repository: `git tag v1.0.4`.
6. Push the tag: `git push origin main v1.0.4`.

## Update Workflow & Verification
1. **Modify:** Update messaging service logic in `app/` or `cmd/`.
2. **Verify Local:**
   - Run `go mod tidy`.
   - Run `go test ./...`.
   - Run `go build -o sb-messenger ./cmd/sb-messenger`.
3. **Commit:** Ensure the commit message clearly describes the messaging change.
4. **Tag & Push:** (Follow the Publishing Order above).
