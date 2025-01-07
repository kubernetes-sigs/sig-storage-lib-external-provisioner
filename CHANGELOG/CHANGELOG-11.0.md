# Release notes for v11.0.0

# Changelog since v10.0.1

## Changes by Kind

### Feature

- Signature of NewProvisionController has been changed to support contextual logging for eventRecorder. This is a breaking change that requires updates to code that calls this function. ([#171](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/171), [@bells17](https://github.com/bells17))

### Bug or Regression

- Add a callback to customize PV name ([#178](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/178), [@smileusd](https://github.com/smileusd))

### Uncategorized

- Fixed overwriting of internal informer cache. This could lead to multiple Provision() calls for a single PVC in very rare (impossible?) cases. The library relies on idempotency of the Provision() call. Please report any frequent duplicate Provision() calls. ([#179](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/179), [@goushicui](https://github.com/goushicui))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._
