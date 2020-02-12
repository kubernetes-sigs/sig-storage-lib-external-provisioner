This directory mirrors the source code via symlinks.
This makes it possible to vendor v5.x releases of
sig-storage-lib-external-provisioner with `dep` versions that do not
support semantic imports. Support for that is currently
[pending in dep](https://github.com/golang/dep/pull/1963).

If users of dep have enabled pruning, they must disable if
for sig-storage-lib-external-provisioner in their Gopk.toml, like this:

```toml
[prune]
  go-tests = true
  unused-packages = true

  [[prune.project]]
    name = "sigs.k8s.io/sig-storage-lib-external-provisioner"
    unused-packages = false
```
