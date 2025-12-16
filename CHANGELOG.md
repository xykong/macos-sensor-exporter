# Changelog

## [1.1.0] - 2025-12-16

### Major Improvements - Code Modernization

This update brings the 3-year-old codebase up to modern Go standards and best practices.

### Fixed

- **Critical Bug**: Fixed `show.go` viper flag binding that was incorrectly binding to `startCmd` instead of `showCmd`
- **Prometheus Compliance**: Implemented proper `Describe()` method using `prometheus.DescribeByCollect()` for dynamic metrics

### Changed

- **Go Version**: Upgraded from Go 1.18 to Go 1.22 (with toolchain 1.24)
- **Dependencies**: Updated all major dependencies to latest compatible versions:
  - `dkorunic/iSMC`: v0.6.2 → v0.7.0
  - `prometheus/client_golang`: v1.12.2 → v1.23.2
  - `spf13/viper`: v1.12.0 → v1.21.0
  - `spf13/cobra`: v1.4.0 → v1.10.2
  - `sirupsen/logrus`: v1.8.1 → v1.9.3
  - Many indirect dependencies also updated for security and compatibility
  
  Note: iSMC v0.10.1 requires Go 1.25+, so we use v0.7.0 which is the latest version compatible with Go 1.22

### Improved

#### Error Handling & Logging
- Replaced `fmt.Print*` with structured logging using `logrus`
- Added proper error context with `WithError()` and `WithField()`
- Improved error messages for better debugging

#### HTTP Server Structure
- Migrated from global `http.Handle()` to explicit `http.NewServeMux()`
- Added `/healthz` health check endpoint
- Proper error handling with `http.ErrServerClosed` check
- Server now exits with `Fatal` on errors instead of just logging

#### Configuration Management
- Enhanced config loading with multiple search paths (home directory + current directory)
- Proper error classification (file not found vs. parse errors)
- Better debug logging for configuration resolution

### Added

- **Comprehensive README**: Added full documentation with:
  - Installation instructions
  - Usage examples
  - Configuration guide
  - Prometheus setup examples
  - Metric format documentation
  - Troubleshooting section
  - Chinese translation (README_zh.md) for better accessibility
  
- **Unit Tests**: Added test coverage for core functionality:
  - `TestGetUnit`: Unit parsing logic
  - `TestGetGaugeValue`: Value conversion logic
  - `TestCreateNewDesc`: Metric descriptor creation
  - `TestSensorsCollectorDescribe`: Prometheus Describe implementation
  - `TestSensorsCollectorCollect`: Full collection workflow
  - All tests passing ✅

### Technical Debt Addressed

1. Removed commented-out code and updated placeholder comments
2. Improved code structure and organization
3. Better separation of concerns in HTTP server setup
4. More idiomatic Go patterns throughout

### Validation

- ✅ All tests passing
- ✅ Build successful
- ✅ No compiler warnings
- ✅ Dependencies up-to-date and compatible

### Future Considerations (Not Implemented)

The following improvements were identified but not implemented in this round:

1. **Metric Naming**: Consider restructuring metric names to follow stricter Prometheus conventions (`namespace_subsystem_name`)
2. **Logging Framework**: Consider migrating from `logrus` to standard library `log/slog` (Go 1.21+)
3. **Testing**: Add integration tests and benchmarks
4. **CI/CD**: Add GitHub Actions for automated testing and releases
5. **Interface Abstraction**: Extract SMC access behind an interface for better testability
