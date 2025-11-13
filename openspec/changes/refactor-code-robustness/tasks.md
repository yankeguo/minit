## 1. Resource Cleanup and Lifecycle Management

- [x] 1.1 Add graceful shutdown support to WebDAV server in `msetups/setup_webdav.go`
  - Accept context for cancellation
  - Call `Server.Shutdown()` on context cancellation
  - Add test for graceful shutdown
- [x] 1.2 Fix timer cleanup in daemon restart logic in `mrunners/runner_daemon.go`
  - Ensure timer is stopped when context is cancelled
  - Add `defer timer.Stop()` or use timer channel with cleanup
- [x] 1.3 Improve file descriptor management in rotating logs in `mlog/rotating.go`
  - Add defensive nil checks before fd operations
  - Ensure fd.Close() is called even on errors
  - Add tests for rotation error scenarios

## 2. Enhanced Error Context and Validation

- [x] 2.1 Improve YAML parsing error messages in `munit/load_file.go`
  - Wrap errors with file path context
  - Add line number information when available
  - Test with malformed YAML files
- [x] 2.2 Enhance unit name validation in `munit/load.go`
  - Include the invalid name in error messages
  - Provide clear validation rules in error text
  - Add tests for various invalid names
- [ ] 2.3 Add environment variable validation in `msetups/setup_rlimits.go`
  - Validate MINIT_RLIMIT_* format before parsing
  - Provide clear error messages for invalid formats
  - Add tests for invalid rlimit values
- [ ] 2.4 Add environment variable validation in `msetups/setup_sysctl.go`
  - Validate sysctl key-value pairs
  - Handle invalid formats gracefully
  - Add tests for edge cases

## 3. Concurrency Safety Improvements

- [x] 3.1 Review and enhance locking in `mlog/rotating.go`
  - Verify atomic operations are sufficient
  - Add comments explaining synchronization strategy
  - Test concurrent writes during rotation
- [x] 3.2 Review signal broadcast synchronization in `mexec/manager.go`
  - Verify lock coverage is correct
  - Add comments explaining thread-safety guarantees
  - Add race detector tests
- [ ] 3.3 Review error group concurrency in `merrs/errors.go`
  - Verify RWMutex usage is optimal
  - Add tests for concurrent Add/Set operations

## 4. Defensive Programming for Edge Cases

- [ ] 4.1 Add nil checks in `mexec/manager.go`
  - Check charset decoder before use
  - Verify cmd.Process is not nil
  - Add tests for nil scenarios
- [x] 4.2 Add empty command validation in `munit/unit.go`
  - Check Command array length in RequireCommand()
  - Provide clear error message
  - Add tests for empty commands
- [ ] 4.3 Improve directory creation in `mlog/logger.go`
  - Use os.MkdirAll for parent directories
  - Provide clear error on permission failure
  - Add tests for directory creation edge cases
- [ ] 4.4 Add bounds checking in array/slice operations
  - Review all slice access patterns
  - Add defensive checks where needed
  - Test with empty and single-element collections

## 5. Input Validation Strengthening

- [ ] 5.1 Add unit order validation in `munit/load.go`
  - Detect potential issues with ordering
  - Document ordering behavior clearly
  - Add tests for edge case ordering
- [x] 5.2 Enhance cron expression validation in `mrunners/runner_cron.go`
  - Validate before runner creation (already done)
  - Improve error message to include unit name
  - Add tests for invalid cron expressions
- [ ] 5.3 Add shell command validation in `mexec/manager.go`
  - Validate shell command can be split correctly
  - Handle empty shell command gracefully
  - Add tests for invalid shell configurations

## 6. Graceful Shutdown Enhancement

- [x] 6.1 Review main shutdown sequence in `main.go`
  - Document shutdown timing and sequence
  - Consider making delay configurable (future enhancement)
  - Add comments explaining the 3-second delay rationale
- [ ] 6.2 Ensure all goroutines are tracked
  - Review all goroutine launches
  - Verify proper WaitGroup usage
  - Add tests for shutdown completeness

## 7. Enhanced Observability and Logging

- [ ] 7.1 Standardize error wrapping across packages
  - Use consistent error format: "operation: detail: underlying_error"
  - Ensure unit name is included in unit-related errors
  - Review all error return statements
- [x] 7.2 Log non-critical errors appropriately in `main.go`
  - Log pprof server errors to stderr or logger
  - Review other silently ignored errors
  - Add appropriate log levels
- [ ] 7.3 Improve error messages in `menv/construct.go`
  - Include variable name in template errors
  - Provide clearer context for rendering failures
  - Add tests for various template errors

## 8. Test Coverage for Edge Cases

- [ ] 8.1 Add race detector tests to CI/CD
  - Ensure `go test -race` passes
  - Document any known benign races (if any)
- [ ] 8.2 Add error path tests for each package
  - Test with invalid inputs
  - Test with nil/empty values
  - Test with boundary conditions
- [ ] 8.3 Add integration tests for shutdown scenarios
  - Test graceful shutdown with multiple units
  - Test signal propagation
  - Test resource cleanup verification
- [ ] 8.4 Improve test coverage metrics
  - Run coverage analysis: `go test -cover ./...`
  - Identify untested code paths
  - Add tests to reach >80% coverage where practical

## 9. Code Review and Documentation

- [ ] 9.1 Add code comments for complex logic
  - Document concurrency patterns
  - Explain error handling strategies
  - Clarify resource ownership
- [ ] 9.2 Review and update package documentation
  - Ensure each package has clear godoc
  - Document exported functions thoroughly
  - Add examples where helpful
- [ ] 9.3 Update error messages for consistency
  - Use consistent terminology
  - Ensure user-facing messages are helpful
  - Include actionable guidance where possible

## 10. Final Validation

- [x] 10.1 Run full test suite with race detector
  - `go test -race ./...`
  - Fix any issues found
- [x] 10.2 Run static analysis tools
  - `go vet ./...`
  - Consider running `staticcheck` if available
- [ ] 10.3 Manual testing of key scenarios
  - Test with various unit configurations
  - Test graceful shutdown
  - Test error scenarios
- [x] 10.4 Review all changes for breaking changes
  - Ensure no API changes
  - Verify backward compatibility
  - Update CHANGELOG if needed

