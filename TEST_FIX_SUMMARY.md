# GitHub Actions Test Fix Summary

## Problem
GitHub Actions workflow was failing during the `go test -v ./...` step with error:
```
Unable to open Apple SMC; return code 1
FAIL	github.com/xykong/macos-sensor-exporter/exporter	0.008s
```

## Root Cause
The tests `TestSensorsCollectorDescribe` and `TestSensorsCollectorCollect` were attempting to access Apple SMC hardware, which is not available on GitHub Actions runners due to virtualization restrictions. Even macOS runners don't have direct hardware access to SMC.

## Solution
Updated `exporter/exporter_test.go` to skip SMC-dependent tests when:
1. Running on non-macOS systems (`runtime.GOOS != "darwin"`)
2. SMC is not accessible (checked by calling `output.GetAll()`)

### Changes Made

**File: `exporter/exporter_test.go`**
- Added `runtime` and `github.com/xykong/iSMC/output` imports
- Added SMC availability checks to `TestSensorsCollectorDescribe` and `TestSensorsCollectorCollect`
- Tests will be skipped with clear messages when SMC is unavailable

**File: `.github/workflows/release.yml`**
- Added comment documenting that SMC-dependent tests will be skipped on GitHub Actions runners

## Test Results
After the fix:
- ✅ Unit tests for helper functions (`TestGetUnit`, `TestGetGaugeValue`, `TestCreateNewDesc`) pass on all platforms
- ✅ SMC-dependent tests skip gracefully on GitHub Actions runners
- ✅ SMC-dependent tests run successfully on physical macOS machines with SMC access
- ✅ No impact on release workflow

## Benefits
1. **CI/CD Reliability**: Tests no longer fail on GitHub Actions
2. **Flexibility**: Tests still run on physical macOS machines with SMC access
3. **Clear Messaging**: Skip messages clearly explain why tests were skipped
4. **No Functionality Loss**: All non-hardware-dependent tests still run

## Testing
Run tests locally:
```bash
cd macos-sensor-exporter
go test -v ./...
```

On physical macOS machines with SMC access, all tests will run.
On GitHub Actions or VMs without SMC access, SMC-dependent tests will be skipped.
