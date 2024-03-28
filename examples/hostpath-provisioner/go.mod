module sigs.k8s.io/kubernetes-sigs/sig-storage-lib-external-provisioner/examples/hostpath-provisioner

go 1.16

require (
	k8s.io/api v0.28.0
	k8s.io/apimachinery v0.28.0
	k8s.io/client-go v0.28.0
	k8s.io/klog/v2 v2.100.1
	sigs.k8s.io/sig-storage-lib-external-provisioner/v10 v9.0.1
)

replace sigs.k8s.io/sig-storage-lib-external-provisioner/v10 => ../..
