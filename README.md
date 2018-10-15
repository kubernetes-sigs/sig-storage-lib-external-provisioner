# sig-storage-lib-external-provisioner

A library for writing [external provisioners](https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner). See [external-storage](https://github.com/kubernetes-incubator/external-storage) for community-maintained external provisioners that use this library.

## Packages
### `controller`
Contains the Provisioner interface and ProvisionController, a custom Kubernetes [controller](https://github.com/kubernetes/community/blob/master/contributors/devel/controllers.md) that watches PersistentVolumes and PersistentVolumeClaims. Implement the Provisioner interface, pass the implementation to a ProvisionController, and Run the controller, which then takes care of calling the Provisioner's Provision or Delete as needed.

## Optional Packages
### `util`
Contains an assortment of useful functions, e.g. any used by [in-tree plugins](https://github.com/kubernetes/kubernetes/tree/master/pkg/volume) that aren't otherwise easily importable.

### `gidallocator` and `allocator`
`gidallocator` is used to allocate a GID from a range specified by StorageClass parameters gidMin & gidMax. `allocator` is the underlying implementation and can be used to write other allocators. An example use-case for `gidallocator` is an NFS-based provisioner that chowns each export to a unique GID. See [Volume Security](https://docs.openshift.com/container-platform/3.11/install_config/persistent_storage/pod_security_context.html#supplemental-groups/) for more context. Adapted from the in-tree [glusterfs plugin](https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/glusterfs).

### `mount`
Is used to read the mount table. Copied from [moby](https://github.com/moby/moby/tree/17.05.x/pkg/mount).

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack](http://slack.k8s.io/)
- [Mailing List](https://groups.google.com/forum/#!forum/kubernetes-dev)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
