# Changelog

## v0.1.0 - 2025-05-23

### Added
- Token registration and configuration via `credence.Config`
- `Token()` function with automatic token caching and refresh
- In-memory `TokenStore` for managing credential lifecycles
- Basic error handling for unregistered configs and missing tokens
- Unit tests for store behavior and error paths

### Notes
- This is the initial MVP release
- Future releases will include Redis support, retry-aware HTTP client, and token lifecycle diagnostics