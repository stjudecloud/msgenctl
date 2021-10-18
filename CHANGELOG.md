# Changelog

## Unreleased

### Added

  * internal/client: Retry on failed HTTP requests.

## 0.2.0 - 2021-09-15

### Changed

  * cmd/wait: Return a failure exit status (non-zero) upon a a failed or
    cancelled workflow status.

    This previously would always return a success exit status (0) even if the
    workflow failed or was cancelled.

## 0.1.0 - 2021-08-31

  * Initial release.
