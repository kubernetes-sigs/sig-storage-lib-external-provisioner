module sigs.k8s.io/kubernetes-sigs/sig-storage-lib-external-provisioner/examples/hostpath-provisioner

go 1.13

require (
	k8s.io/api v0.19.1
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v0.19.1
	k8s.io/klog/v2 v2.3.0
	sigs.k8s.io/sig-storage-lib-external-provisioner/v6 v6.0.0
)

replace sigs.k8s.io/sig-storage-lib-external-provisioner/v6 => ../..
