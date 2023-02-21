module sigs.k8s.io/kubernetes-sigs/sig-storage-lib-external-provisioner/examples/hostpath-provisioner

go 1.16

require (
	k8s.io/api v0.19.1
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v0.19.1
	k8s.io/klog/v2 v2.3.0
	sigs.k8s.io/sig-storage-lib-external-provisioner/v9 v9.0.1
)

replace sigs.k8s.io/sig-storage-lib-external-provisioner/v9 => ../..
