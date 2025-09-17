# Release notes for v13.0.0

# Changelog since v12.0.1

## Changes by Kind

### Urgent Upgrade Notes

- Update to v13: ProvisionOptions now gets just name of the selected node instead of the whole *v1.Node object. NodeLister has been removed from the ProvisionController. Individual provisioners need to fetch the node on their own, if they need it. ([#194](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/194), [@sunnylovestiramisu](https://github.com/sunnylovestiramisu))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._
