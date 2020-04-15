module sigs.k8s.io/kubernetes-sigs/sig-storage-lib-external-provisioner/examples/hostpath-provisioner

go 1.13

require (
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	k8s.io/klog v1.0.0
	sigs.k8s.io/sig-storage-lib-external-provisioner/v5 v5.0.0
)

replace sigs.k8s.io/sig-storage-lib-external-provisioner/v5 => ../..
