module sigs.k8s.io/kubernetes-sigs/sig-storage-lib-external-provisioner/examples/hostpath-provisioner

go 1.22.0

toolchain go1.22.2

require (
	k8s.io/api v0.30.0
	k8s.io/apimachinery v0.30.0
	k8s.io/client-go v0.30.0
	k8s.io/klog/v2 v2.120.1
	sigs.k8s.io/sig-storage-lib-external-provisioner/v10 v10.0.0-20240423100449-ea3e5f96b47e
)

replace sigs.k8s.io/sig-storage-lib-external-provisioner/v10 => ../..
