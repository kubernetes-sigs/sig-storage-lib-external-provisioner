/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	v1 "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	testclient "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	ref "k8s.io/client-go/tools/reference"
	"k8s.io/client-go/util/workqueue"
	klog "k8s.io/klog/v2"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v9/controller/metrics"
)

const (
	resyncPeriod         = 100 * time.Millisecond
	sharedResyncPeriod   = 1 * time.Second
	defaultServerVersion = "v1.5.0"
)

func init() {
	klog.InitFlags(nil)
}

var (
	modeWait = storage.VolumeBindingWaitForFirstConsumer
)

// TODO clean this up, e.g. remove redundant params (provisionerName: "foo.bar/baz")
func TestController(t *testing.T) {
	var reflectorCallCount int
	timestamp := metav1.NewTime(time.Now())
	tests := []struct {
		name                       string
		objs                       []runtime.Object
		claimsInProgress           []*v1.PersistentVolumeClaim
		enqueueClaim               *v1.PersistentVolumeClaim
		provisionerName            string
		additionalProvisionerNames []string
		provisioner                Provisioner
		verbs                      []string
		reaction                   testclient.ReactionFunc
		expectedVolumes            []v1.PersistentVolume
		expectedClaims             []v1.PersistentVolumeClaim
		expectedClaimsInProgress   []string
		volumeQueueStore           bool
		expectedStoredVolumes      []*v1.PersistentVolume
		expectedMetrics            testMetrics
		deletionTimestamp          *metav1.Time
		addFinalizer               bool
	}{
		{
			name: "provision for claim-1 but not claim-2",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newStorageClass("class-2", "abc.def/ghi"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
				newClaim("claim-2", "uid-1-2", "class-2", "abc.def/ghi", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "don't provision, volume already exists",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
				newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "don't provision, provisioner does not support raw block",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaimWithVolumeMode("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil, v1.PersistentVolumeBlock),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "provision for claim-1 with storage class provisioner name distinct from controller provisioner name",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName:            "csi.com/mock-csi",
			additionalProvisionerNames: []string{"foo.bar/baz", "foo.xyz/baz"},
			provisioner:                newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "delete volume-1 but not volume-2",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
				newVolume("volume-2", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "abc.def/ghi"}, nil, nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-2", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "abc.def/ghi"}, nil, nil),
			},
			expectedMetrics: testMetrics{
				deleted: counts{
					"": count{success: 1},
				},
			},
		},
		{
			name: "don't provision for claim-1 because it's already bound",
			objs: []runtime.Object{
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "volume-1", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume(nil),
		},
		{
			name: "don't provision for claim-1 because its class doesn't exist",
			objs: []runtime.Object{
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume(nil),
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "don't delete volume-1 because it's still bound",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
		},
		{
			name: "don't delete volume-1 because its reclaim policy is not delete",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimRetain, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimRetain, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
		},
		{
			name: "provisioner fails to provision for claim-1: no pv is created",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newBadTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume(nil),
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "provisioner fails to delete volume-1: pv is not deleted",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newBadTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
			expectedMetrics: testMetrics{
				deleted: counts{
					"": count{failed: 1},
				},
			},
		},
		{
			name: "try to provision for claim-1 but fail to save the pv object",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			verbs:           []string{"create"},
			reaction: func(action testclient.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, errors.New("fake error")
			},
			expectedVolumes: []v1.PersistentVolume(nil),
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "try to delete volume-1 but fail to delete the pv object",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			verbs:           []string{"delete"},
			reaction: func(action testclient.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, errors.New("fake error")
			},
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
			expectedMetrics: testMetrics{
				deleted: counts{
					"": count{failed: 1},
				},
			},
		},
		{
			name: "don't provision, because it is ignored",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-2", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newIgnoredProvisioner(),
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "provision for claim-1 but not claim-2, because it is ignored",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
				newClaim("claim-2", "uid-1-2", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newIgnoredProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1, failed: 1},
				},
			},
		},
		{
			name: "provision with Retain reclaim policy",
			objs: []runtime.Object{
				newStorageClassWithReclaimPolicy("class-1", "foo.bar/baz", v1.PersistentVolumeReclaimRetain),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolumeWithReclaimPolicy(newStorageClassWithReclaimPolicy("class-1", "foo.bar/baz", v1.PersistentVolumeReclaimRetain), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "provision for ext provisioner",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newProvisioner(t, "pvc-uid-1-1", ProvisioningFinished, nil),
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "ext provisioner: final error does not mark claim as in progress",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName:          "foo.bar/baz",
			provisioner:              newProvisioner(t, "pvc-uid-1-1", ProvisioningFinished, fmt.Errorf("mock error")),
			expectedClaimsInProgress: []string{},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "ext provisioner: provisional error marks claim as in progress",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName:          "foo.bar/baz",
			provisioner:              newProvisioner(t, "pvc-uid-1-1", ProvisioningInBackground, fmt.Errorf("mock error")),
			expectedClaimsInProgress: []string{"uid-1-1"},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "ext provisioner: NoChange error does not mark claim as in progress",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName:          "foo.bar/baz",
			provisioner:              newProvisioner(t, "pvc-uid-1-1", ProvisioningNoChange, fmt.Errorf("mock error")),
			expectedClaimsInProgress: []string{},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "ext provisioner: final error removes claim from in progress",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			claimsInProgress: []*v1.PersistentVolumeClaim{
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName:          "foo.bar/baz",
			provisioner:              newProvisioner(t, "pvc-uid-1-1", ProvisioningFinished, fmt.Errorf("mock error")),
			expectedClaimsInProgress: []string{},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "ext provisioner: provisional error does not remove claim from in progress",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			claimsInProgress: []*v1.PersistentVolumeClaim{
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName:          "foo.bar/baz",
			provisioner:              newProvisioner(t, "pvc-uid-1-1", ProvisioningInBackground, fmt.Errorf("mock error")),
			expectedClaimsInProgress: []string{"uid-1-1"},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "ext provisioner: NoChange error does not remove claim from in progress",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			claimsInProgress: []*v1.PersistentVolumeClaim{
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName:          "foo.bar/baz",
			provisioner:              newProvisioner(t, "pvc-uid-1-1", ProvisioningNoChange, fmt.Errorf("mock error")),
			expectedClaimsInProgress: []string{"uid-1-1"},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "ext provisioner: claimsInProgress is used for deleted PVCs",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
			},
			claimsInProgress: []*v1.PersistentVolumeClaim{
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			enqueueClaim:             newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			provisionerName:          "foo.bar/baz",
			provisioner:              newProvisioner(t, "pvc-uid-1-1", ProvisioningFinished, nil),
			expectedClaimsInProgress: []string{},
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "PV save backoff: provision a PV and fail to save it -> it's in the queue",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			verbs:           []string{"create"},
			reaction: func(action testclient.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, errors.New("fake error")
			},
			expectedVolumes:  []v1.PersistentVolume(nil),
			volumeQueueStore: true,
			expectedStoredVolumes: []*v1.PersistentVolume{
				newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "PV save backoff: provision a PV and fail to save it two times -> it's removed from the queue",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			verbs:           []string{"create"},
			reaction: func(action testclient.Action) (handled bool, ret runtime.Object, err error) {
				reflectorCallCount++
				if reflectorCallCount <= 2 {
					return true, nil, errors.New("fake error")
				}
				return false, nil, nil
			},
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
			volumeQueueStore:      true,
			expectedStoredVolumes: []*v1.PersistentVolume{},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "remove selectedNode and claim on reschedule",
			objs: []runtime.Object{
				newStorageClassWithVolumeBindingMode("class-1", "foo.bar/baz", &modeWait),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node-1"}),
				newNode("node-1"),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newRescheduleTestProvisioner(),
			expectedClaims: []v1.PersistentVolumeClaim{
				*newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz"}),
			},
			expectedClaimsInProgress: nil, // not in progress anymore
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "do not remove selectedNode after final error, only the claim",
			objs: []runtime.Object{
				newStorageClassWithVolumeBindingMode("class-1", "foo.bar/baz", &modeWait),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node-wrong"}),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newBadTestProvisioner(),
			expectedClaims: []v1.PersistentVolumeClaim{
				*newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node-wrong"}),
			},
			expectedClaimsInProgress: nil, // not in progress anymore
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "do not remove selectedNode if nothing changes",
			objs: []runtime.Object{
				newStorageClassWithVolumeBindingMode("class-1", "foo.bar/baz", &modeWait),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node-wrong"}),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newNoChangeTestProvisioner(),
			expectedClaims: []v1.PersistentVolumeClaim{
				*newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node-wrong"}),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "do not remove selectedNode while in progress",
			objs: []runtime.Object{
				newStorageClassWithVolumeBindingMode("class-1", "foo.bar/baz", &modeWait),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node-1"}),
				newNode("node-1"),
			},
			provisionerName: "foo.bar/baz",
			provisioner:     newTemporaryTestProvisioner(),
			expectedClaims: []v1.PersistentVolumeClaim{
				*newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node-1"}),
			},
			expectedClaimsInProgress: []string{"uid-1-1"},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{failed: 1},
				},
			},
		},
		{
			name: "provision for ext provisioner with finalizer",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newProvisioner(t, "pvc-uid-1-1", ProvisioningFinished, nil),
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), []string{finalizerPV}),
			},
			expectedMetrics: testMetrics{
				provisioned: counts{
					"class-1": count{success: 1},
				},
			},
		},
		{
			name: "ensure finalizer is removed if the addFinalizer config option is false",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
			addFinalizer:    false,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
		},
		{
			name: "ensure finalizer is removed if the reclaim policy is Retain or Recycle with addFinalizer enabled",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimRetain, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimRetain, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
		},
		{
			name: "ensure finalizer is not added if the volume is under deletion, also ensures that volume is not deleted",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, &timestamp),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, &timestamp),
			},
		},
		{
			name: "ensure volume with finalizer is deleted if it is in a Released state and reclaim policy is Delete",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedMetrics: testMetrics{
				deleted: counts{
					"": count{success: 1},
				},
			},
		},
		{
			name: "ensure volume with finalizer is deleted if it is in a Released state and reclaim policy is Delete and already under deletion",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, &timestamp),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedMetrics: testMetrics{
				deleted: counts{
					"": count{success: 1},
				},
			},
		},
		{
			name: "provisioner fails to delete the volume with finalizer, the pv is not deleted",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newBadTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
			expectedMetrics: testMetrics{
				deleted: counts{
					"": count{failed: 1},
				},
			},
		},
		{
			name: "volume deletion succeeds but the pv deletion fails, the pv still exists",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			verbs:           []string{"delete"},
			reaction: func(action testclient.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, errors.New("fake error")
			},
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
			expectedMetrics: testMetrics{
				deleted: counts{
					"": count{failed: 1},
				},
			},
		},
		{
			name: "ensure finalizer is added on statically provisioned, migrated, in-tree volumes",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annMigratedTo: "foo.bar/baz"}, nil, nil),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annMigratedTo: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
		},
		{
			name: "ensure finalizer is added on statically provisioned, migrated, in-tree volumes if it is in a Bound state",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annMigratedTo: "foo.bar/baz"}, nil, nil),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annMigratedTo: "foo.bar/baz"}, []string{finalizerPV}, nil),
			},
		},
		{
			name: "ensure finalizer is added on statically provisioned CSI volumes if it is in a Bound state",
			objs: []runtime.Object{
				newCSIVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, nil, nil, nil, "foo.bar/baz"),
			},
			addFinalizer:    true,
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			expectedVolumes: []v1.PersistentVolume{
				*newCSIVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, nil, []string{finalizerPV}, nil, "foo.bar/baz"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reflectorCallCount = 0

			client := fake.NewSimpleClientset(test.objs...)
			if len(test.verbs) != 0 {
				for _, v := range test.verbs {
					client.Fake.PrependReactor(v, "persistentvolumes", test.reaction)
				}
			}

			var ctrl testProvisionController
			provisionerOptions := make([]func(*ProvisionController) error, 0)
			if test.addFinalizer {
				provisionerOptions = append(provisionerOptions, AddFinalizer(true))
			}
			if test.additionalProvisionerNames == nil {
				ctrl = newTestProvisionController(client, test.provisionerName, test.provisioner, provisionerOptions...)
			} else {
				ctrl = newTestProvisionControllerWithAdditionalNames(client, test.provisionerName, test.provisioner, test.additionalProvisionerNames, provisionerOptions...)
			}
			for _, claim := range test.claimsInProgress {
				ctrl.claimsInProgress.Store(string(claim.UID), claim)
			}

			if test.volumeQueueStore {
				ctrl.volumeStore = NewVolumeStoreQueue(client, workqueue.DefaultItemBasedRateLimiter(), ctrl.claimsIndexer, ctrl.eventRecorder)
			}

			if test.enqueueClaim != nil {
				ctrl.enqueueClaim(test.enqueueClaim)
			}

			// Run forever...
			go ctrl.Run(context.Background())

			// When we shutdown while something is happening the fake client panics
			// with send on closed channel...but the test passed, so ignore
			utilruntime.ReallyCrash = false

			time.Sleep(2 * resyncPeriod)

			pvList, _ := client.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
			if !reflect.DeepEqual(test.expectedVolumes, pvList.Items) {
				t.Errorf("expected PVs:\n %v\n but got:\n %v\n", test.expectedVolumes, pvList.Items)
			}

			claimsInProgress := sets.NewString()
			ctrl.claimsInProgress.Range(func(key, value interface{}) bool {
				claimsInProgress.Insert(key.(string))
				return true
			})
			expectedClaimsInProgress := sets.NewString(test.expectedClaimsInProgress...)
			if !claimsInProgress.Equal(expectedClaimsInProgress) {
				t.Errorf("expected claimsInProgres: %+v, got %+v", expectedClaimsInProgress.List(), claimsInProgress.List())
			}

			if test.volumeQueueStore {
				queue := ctrl.volumeStore.(*queueStore)
				// convert queue.volumes to array
				queuedVolumes := []*v1.PersistentVolume{}
				queue.volumes.Range(func(key, value interface{}) bool {
					volume, ok := value.(*v1.PersistentVolume)
					if !ok {
						t.Errorf("Expected PersistentVolume in volume store queue, got %+v", value)
					}
					queuedVolumes = append(queuedVolumes, volume)
					return true
				})
				if !reflect.DeepEqual(test.expectedStoredVolumes, queuedVolumes) {
					t.Errorf("expected stored volumes:\n %v\n got: \n%v", test.expectedStoredVolumes, queuedVolumes)
				}

				// Check that every volume is really in the workqueue. It has no List() functionality, use NumRequeues
				// as workaround.
				for _, volume := range test.expectedStoredVolumes {
					if queue.queue.NumRequeues(volume.Name) == 0 {
						t.Errorf("Expected volume %q in workqueue, but it has zero NumRequeues", volume.Name)
					}
				}
			}

			tm := ctrl.getMetrics(t)
			if !reflect.DeepEqual(test.expectedMetrics, tm) {
				t.Errorf("expected metrics:\n %+v\n but got:\n %+v", test.expectedMetrics, tm)
			}

			if test.expectedClaims != nil {
				pvcList, _ := client.CoreV1().PersistentVolumeClaims(v1.NamespaceDefault).List(context.Background(), metav1.ListOptions{})
				if !reflect.DeepEqual(test.expectedClaims, pvcList.Items) {
					t.Errorf("expected PVCs:\n %v\n but got:\n %v\n", test.expectedClaims, pvcList.Items)
				}
			}
		})
	}
}

func TestTopologyParams(t *testing.T) {
	dummyAllowedTopology := []v1.TopologySelectorTerm{
		{
			MatchLabelExpressions: []v1.TopologySelectorLabelRequirement{
				{
					Key:    "failure-domain.beta.kubernetes.io/zone",
					Values: []string{"zone1"},
				},
			},
		},
	}

	tests := []struct {
		name           string
		objs           []runtime.Object
		expectedParams *provisionParams
	}{
		{
			name: "provision without topology information",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			expectedParams: &provisionParams{},
		},
		{
			name: "provision with AllowedTopologies",
			objs: []runtime.Object{
				newStorageClassWithAllowedTopologies("class-1", "foo.bar/baz", dummyAllowedTopology),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			expectedParams: &provisionParams{
				allowedTopologies: dummyAllowedTopology,
			},
		},
		{
			name: "provision with selected node",
			objs: []runtime.Object{
				newNode("node-1"),
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annSelectedNode: "node-1"}),
			},
			expectedParams: &provisionParams{
				selectedNode: newNode("node-1"),
			},
		},
		{
			name: "provision with AllowedTopologies and selected node",
			objs: []runtime.Object{
				newNode("node-1"),
				newStorageClassWithAllowedTopologies("class-1", "foo.bar/baz", dummyAllowedTopology),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annSelectedNode: "node-1"}),
			},
			expectedParams: &provisionParams{
				allowedTopologies: dummyAllowedTopology,
				selectedNode:      newNode("node-1"),
			},
		},
		{
			name: "provision with selected node, but node does not exist",
			objs: []runtime.Object{
				newStorageClass("class-1", "foo.bar/baz"),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", map[string]string{annSelectedNode: "node-1"}),
			},
			expectedParams: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.objs...)
			provisioner := newTestProvisioner()
			ctrl := newTestProvisionController(client, "foo.bar/baz" /* provisionerName */, provisioner)
			// Run forever...
			go ctrl.Run(context.Background())

			// When we shutdown while something is happening the fake client panics
			// with send on closed channel...but the test passed, so ignore
			utilruntime.ReallyCrash = false

			time.Sleep(2 * resyncPeriod)

			if test.expectedParams == nil {
				if len(provisioner.provisionCalls) != 0 {
					t.Errorf("did not expect a Provision() call but got at least 1")
				}
			} else {
				if len(provisioner.provisionCalls) == 0 {
					t.Errorf("expected Provision() call but got none")
				} else {
					actual := <-provisioner.provisionCalls
					if !reflect.DeepEqual(*test.expectedParams, actual) {
						t.Errorf("expected topology parameters: %v; actual: %v", test.expectedParams, actual)
					}
				}
			}
		})
	}
}

func TestShouldProvision(t *testing.T) {
	tests := []struct {
		name                       string
		provisionerName            string
		additionalProvisionerNames []string
		provisioner                Provisioner
		class                      *storage.StorageClass
		claim                      *v1.PersistentVolumeClaim
		expectedShould             bool
		expectedError              bool
	}{
		{
			name:            "should provision based on provisionerName",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim:           newClaim("claim-1", "1-1", "class-1", "foo.bar/baz", "", nil),
			expectedShould:  true,
		},
		{
			name:                       "should provision based on additionalProvisionerNames",
			provisionerName:            "csi.com/mock-csi",
			additionalProvisionerNames: []string{"foo.bar/baz", "foo.xyz/baz"},
			provisioner:                newTestProvisioner(),
			class:                      newStorageClass("class-1", "foo.bar/baz"),
			claim:                      newClaim("claim-1", "1-1", "class-1", "foo.bar/baz", "", nil),
			expectedShould:             true,
		},
		{
			name:            "claim already bound",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim:           newClaim("claim-1", "1-1", "class-1", "foo.bar/baz", "foo", nil),
			expectedShould:  false,
		},
		{
			name:            "no such class",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim:           newClaim("claim-1", "1-1", "class-2", "", "", nil),
			expectedShould:  false,
		},
		{
			name:            "not this provisioner's job",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClass("class-1", "abc.def/ghi"),
			claim:           newClaim("claim-1", "1-1", "class-1", "abc.def/ghi", "", nil),
			expectedShould:  false,
		},
		// Kubernetes 1.5 provisioning - annBetaStorageProvisioner is set
		// and only this annotation is evaluated
		{
			name:            "unknown provisioner annotation 1.5",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim: newClaim("claim-1", "1-1", "class-1", "", "",
				map[string]string{annBetaStorageProvisioner: "abc.def/ghi"}),
			expectedShould: false,
		},
		// Kubernetes 1.5 provisioning - annBetaStorageProvisioner is not set
		{
			name:            "no provisioner annotation 1.5",
			provisionerName: "foo.bar/baz",
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim:           newClaim("claim-1", "1-1", "class-1", "", "", nil),
			expectedShould:  false,
		},
		// Kubernetes 1.23 provisioning - annStorageProvisioner is set
		{
			name:            "unknown provisioner annotation 1.23",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim: newClaim("claim-1", "1-1", "class-1", "", "",
				map[string]string{annStorageProvisioner: "abc.def/ghi"}),
			expectedShould: false,
		},
		{
			name:            "qualifier says no",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestQualifiedProvisioner(false),
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim:           newClaim("claim-1", "1-1", "class-1", "foo.bar/baz", "", nil),
			expectedShould:  false,
		},
		{
			name:            "qualifier says yes, should provision",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestQualifiedProvisioner(true),
			class:           newStorageClass("class-1", "foo.bar/baz"),
			claim:           newClaim("claim-1", "1-1", "class-1", "foo.bar/baz", "", nil),
			expectedShould:  true,
		},
		{
			name:            "if PVC is in delay binding mode, should not provision if annSelectedNode is not set",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClassWithVolumeBindingMode("class-1", "foo.bar/baz", &modeWait),
			claim:           newClaim("claim-1", "1-1", "class-1", "", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz"}),
			expectedShould:  false,
		},
		{
			name:            "if PVC is in delay binding mode, should provision if annSelectedNode is set",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClassWithVolumeBindingMode("class-1", "foo.bar/baz", &modeWait),
			claim:           newClaim("claim-1", "1-1", "class-1", "", "", map[string]string{annBetaStorageProvisioner: "foo.bar/baz", annSelectedNode: "node1"}),
			expectedShould:  true,
		},
		{
			name:            "if PVC is in delay binding mode, should provision if annSelectedNode is set with annStorageProvisioner",
			provisionerName: "foo.bar/baz",
			provisioner:     newTestProvisioner(),
			class:           newStorageClassWithVolumeBindingMode("class-1", "foo.bar/baz", &modeWait),
			claim:           newClaim("claim-1", "1-1", "class-1", "", "", map[string]string{annStorageProvisioner: "foo.bar/baz", annSelectedNode: "node1"}),
			expectedShould:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.claim)

			var ctrl testProvisionController
			if test.additionalProvisionerNames == nil {
				ctrl = newTestProvisionController(client, test.provisionerName, test.provisioner)
			} else {
				ctrl = newTestProvisionControllerWithAdditionalNames(client, test.provisionerName, test.provisioner, test.additionalProvisionerNames)
			}

			if test.class != nil {
				err := ctrl.classes.Add(test.class)
				if err != nil {
					t.Errorf("error adding class %v to cache: %v", test.class, err)
				}
			}

			should, err := ctrl.shouldProvision(context.Background(), test.claim)
			if test.expectedShould != should {
				t.Errorf("expected should provision %v but got %v\n", test.expectedShould, should)
			}
			if (err != nil && test.expectedError == false) || (err == nil && test.expectedError == true) {
				t.Errorf("expected error %v but got %v\n", test.expectedError, err)
			}
		})
	}
}

func TestShouldDelete(t *testing.T) {
	timestamp := metav1.NewTime(time.Now())
	tests := []struct {
		name              string
		provisionerName   string
		volume            *v1.PersistentVolume
		deletionTimestamp *metav1.Time
		expectedShould    bool
	}{
		{
			name:            "should delete",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			expectedShould:  true,
		},
		{
			name:            "failed: shouldn't delete",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeFailed, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			expectedShould:  false,
		},
		{
			name:            "volume still bound",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			expectedShould:  false,
		},
		{
			name:            "non-delete reclaim policy",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimRetain, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			expectedShould:  false,
		},
		{
			name:              "non-nil deletion timestamp",
			provisionerName:   "foo.bar/baz",
			volume:            newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			deletionTimestamp: &timestamp,
			expectedShould:    false,
		},
		{
			name:            "nil deletion timestamp",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			expectedShould:  true,
		},
		{
			name:            "migrated to",
			provisionerName: "csi.driver",
			volume:          newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz", annMigratedTo: "csi.driver"}, nil, nil),
			expectedShould:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := fake.NewSimpleClientset()
			provisioner := newTestProvisioner()
			ctrl := newTestProvisionController(client, test.provisionerName, provisioner)
			test.volume.ObjectMeta.DeletionTimestamp = test.deletionTimestamp

			should := ctrl.shouldDelete(context.Background(), test.volume)
			if test.expectedShould != should {
				t.Errorf("expected should delete %v but got %v\n", test.expectedShould, should)
			}
		})
	}
}

func TestIsProvisionerForVolume(t *testing.T) {
	tests := []struct {
		name            string
		provisionerName string
		volume          *v1.PersistentVolume
		expectedShould  bool
	}{
		{
			name:            "known dynamically provisioned in-tree volume",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeAvailable, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			expectedShould:  true,
		},
		{
			name:            "known dynamically provisioned in-tree migrated volume",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeAvailable, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz", annMigratedTo: "foo.bar/baz"}, nil, nil),
			expectedShould:  true,
		},
		{
			name:            "unknown dynamically provisioned in-tree volume",
			provisionerName: "foo.bar1/baz1",
			volume:          newVolume("volume-1", v1.VolumeAvailable, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			expectedShould:  false,
		},
		{
			name:            "unknown dynamically provisioned in-tree migrated volume",
			provisionerName: "foo.bar1/baz1",
			volume:          newVolume("volume-1", v1.VolumeAvailable, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz", annMigratedTo: "foo.bar/baz"}, nil, nil),
			expectedShould:  false,
		},
		{
			name:            "known dynamically provisioned csi volume",
			provisionerName: "foo.bar/baz",
			volume:          newCSIVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil, "foo.bar/baz"),
			expectedShould:  true,
		},
		{
			name:            "unknown dynamically provisioned csi volume",
			provisionerName: "foo.bar1/baz1",
			volume:          newCSIVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil, "foo.bar/baz"),
			expectedShould:  false,
		},
		{
			name:            "known statically provisioned csi volume",
			provisionerName: "foo.bar/baz",
			volume:          newCSIVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, nil, nil, nil, "foo.bar/baz"),
			expectedShould:  true,
		},
		{
			name:            "unknown statically provisioned csi volume",
			provisionerName: "foo.bar1/baz1",
			volume:          newCSIVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, nil, nil, nil, "foo.bar/baz"),
			expectedShould:  false,
		},
		{
			name:            "statically provisioned in-tree volume always return false", // always return false for statically provisioned in-tree volumes
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeAvailable, v1.PersistentVolumeReclaimDelete, nil, nil, nil),
			expectedShould:  false,
		},
		{
			name:            "known statically provisioned in-tree migrate volume",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeAvailable, v1.PersistentVolumeReclaimDelete, map[string]string{annMigratedTo: "foo.bar/baz"}, nil, nil),
			expectedShould:  true,
		},
		{
			name:            "unknown statically provisioned in-tree migrate volume",
			provisionerName: "foo.bar1/baz1",
			volume:          newVolume("volume-1", v1.VolumeAvailable, v1.PersistentVolumeReclaimDelete, map[string]string{annMigratedTo: "foo.bar/baz"}, nil, nil),
			expectedShould:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := fake.NewSimpleClientset()
			provisioner := newTestProvisioner()
			ctrl := newTestProvisionController(client, test.provisionerName, provisioner)
			should := ctrl.isProvisionerForVolume(context.Background(), test.volume)
			if test.expectedShould != should {
				t.Errorf("expected should delete %v but got %v\n", test.expectedShould, should)
			}
		})
	}
}

func TestShouldDeleteWithFinalizer(t *testing.T) {
	timestamp := metav1.NewTime(time.Now())
	tests := []struct {
		name              string
		provisionerName   string
		volume            *v1.PersistentVolume
		deletionTimestamp *metav1.Time
		expectedShould    bool
	}{
		{
			// Represents a normal deletion where the PVC is deleted and the PV with `Delete` reclaim policy is in
			// `Released` state. The PV has the finalizer that was previously added to it by `syncVolume`.
			name:            "should delete with finalizer present and no deletionTimestamp",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			expectedShould:  true,
		},
		{
			// Represents a scenario where the `syncVolume` has manged to add the finalizer, this is followed by PV
			// deletion first, then the PVC deletion, in this case, assuming the PV has `Delete` reclaim policy, the PV
			// would be in a `Released` state with deletionTimestamp set.
			name:              "should delete with finalizer present and with deletionTimestamp",
			provisionerName:   "foo.bar/baz",
			volume:            newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			deletionTimestamp: &timestamp,
			expectedShould:    true,
		},
		{
			// Represents some update on the PV, the PV has the finalizer due to `syncVolume`. For example, annotation
			// addition or deletion.
			name:            "should not delete when volume still bound and the finalizer exists with no deletionTimestamp",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{finalizerPV}, nil),
			expectedShould:  false,
		},
		{
			// Represents scenario where the reclaim policy is changed to `Retain` or `Recycle`, the syncVolume removes
			// the finalizer. Volume should not be deleted in this case, here, the deletionTimestamp does not matter.
			name:            "should not delete for non-delete reclaim policy",
			provisionerName: "foo.bar/baz",
			volume:          newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimRetain, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, []string{}, nil),
			expectedShould:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := fake.NewSimpleClientset()
			provisioner := newTestProvisioner()
			provisionerOptions := make([]func(*ProvisionController) error, 0)
			provisionerOptions = append(provisionerOptions, AddFinalizer(true))
			ctrl := newTestProvisionController(client, test.provisionerName, provisioner, provisionerOptions...)
			test.volume.ObjectMeta.DeletionTimestamp = test.deletionTimestamp

			should := ctrl.shouldDelete(context.Background(), test.volume)
			if test.expectedShould != should {
				t.Errorf("expected should delete %v but got %v\n", test.expectedShould, should)
			}
		})
	}
}

func TestRemoveFinalizer(t *testing.T) {
	tests := []struct {
		name               string
		finalizers         []string
		finalizerToRemove  string
		expectedFinalizers []string
		expectedModified   bool
	}{
		{
			name:               "ensure the finalizer removal returns nil finalizers and modified",
			finalizers:         []string{finalizerPV},
			finalizerToRemove:  finalizerPV,
			expectedFinalizers: nil,
			expectedModified:   true,
		},
		{
			name:               "if finalizer does not exist then it should not be modified",
			finalizers:         []string{"dummy"},
			finalizerToRemove:  finalizerPV,
			expectedFinalizers: []string{"dummy"},
			expectedModified:   false,
		},
		{
			name:               "remove only the finalizerToAdd from multiple finalizers",
			finalizers:         []string{"dummy", finalizerPV},
			finalizerToRemove:  finalizerPV,
			expectedFinalizers: []string{"dummy"},
			expectedModified:   true,
		},
		{
			name:               "no modified if the original finalizer is empty",
			finalizers:         []string{},
			finalizerToRemove:  finalizerPV,
			expectedFinalizers: nil,
			expectedModified:   false,
		},
		{
			name:               "remove from nil finalizer",
			finalizers:         nil,
			finalizerToRemove:  finalizerPV,
			expectedFinalizers: nil,
			expectedModified:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			modifiedFinalizers, modified := removeFinalizer(test.finalizers, test.finalizerToRemove)
			if test.expectedModified != modified {
				t.Errorf("expected modified %v but got %v\n", test.expectedModified, modified)
			}
			if !reflect.DeepEqual(test.expectedFinalizers, modifiedFinalizers) {
				t.Errorf("expected finalizers %v but got %v\n", test.expectedFinalizers, modifiedFinalizers)
			}
		})
	}
}

func TestAddFinalizer(t *testing.T) {
	tests := []struct {
		name               string
		finalizers         []string
		finalizerToAdd     string
		expectedFinalizers []string
		expectedModified   bool
	}{
		{
			name:               "if finalizer already exists then no modification needed",
			finalizers:         []string{finalizerPV},
			finalizerToAdd:     finalizerPV,
			expectedFinalizers: []string{finalizerPV},
			expectedModified:   false,
		},
		{
			name:               "if original finalizer is nil then return with finalizer to add",
			finalizers:         nil,
			finalizerToAdd:     finalizerPV,
			expectedFinalizers: []string{finalizerPV},
			expectedModified:   true,
		},
		{
			name:               "if original finalizer is empty then return with finalizer to add",
			finalizers:         []string{},
			finalizerToAdd:     finalizerPV,
			expectedFinalizers: []string{finalizerPV},
			expectedModified:   true,
		},
		{
			name:               "ensure finalizers are appended to existing finalizer",
			finalizers:         []string{"dummy"},
			finalizerToAdd:     finalizerPV,
			expectedFinalizers: []string{"dummy", finalizerPV},
			expectedModified:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modifiedFinalizers, modified := addFinalizer(test.finalizers, test.finalizerToAdd)
			if test.expectedModified != modified {
				t.Errorf("expected modified %v but got %v\n", test.expectedModified, modified)
			}
			if !reflect.DeepEqual(test.expectedFinalizers, modifiedFinalizers) {
				t.Errorf("expected finalizers %v but got %v\n", test.expectedFinalizers, modifiedFinalizers)
			}
		})
	}
}

func TestCanProvision(t *testing.T) {
	const (
		provisionerName = "foo.bar/baz"
		blockErrFormat  = "%s does not support block volume provisioning"
	)

	tests := []struct {
		name        string
		provisioner Provisioner
		claim       *v1.PersistentVolumeClaim
		expectedCan error
	}{
		// volumeMode tests for provisioner w/o BlockProvisoner I/F
		{
			name:        "Undefined volumeMode PV request to provisioner w/o BlockProvisoner I/F",
			provisioner: newTestProvisioner(),
			claim:       newClaim("claim-1", "1-1", "class-1", provisionerName, "", nil),
			expectedCan: nil,
		},
		{
			name:        "FileSystem volumeMode PV request to provisioner w/o BlockProvisoner I/F",
			provisioner: newTestProvisioner(),
			claim:       newClaimWithVolumeMode("claim-1", "1-1", "class-1", provisionerName, "", nil, v1.PersistentVolumeFilesystem),
			expectedCan: nil,
		},
		{
			name:        "Block volumeMode PV request to provisioner w/o BlockProvisoner I/F",
			provisioner: newTestProvisioner(),
			claim:       newClaimWithVolumeMode("claim-1", "1-1", "class-1", provisionerName, "", nil, v1.PersistentVolumeBlock),
			expectedCan: fmt.Errorf(blockErrFormat, provisionerName),
		},
		// volumeMode tests for BlockProvisioner that returns false
		{
			name:        "Undefined volumeMode PV request to BlockProvisoner that returns false",
			provisioner: newTestBlockProvisioner(false),
			claim:       newClaim("claim-1", "1-1", "class-1", provisionerName, "", nil),
			expectedCan: nil,
		},
		{
			name:        "FileSystem volumeMode PV request to BlockProvisoner that returns false",
			provisioner: newTestBlockProvisioner(false),
			claim:       newClaimWithVolumeMode("claim-1", "1-1", "class-1", provisionerName, "", nil, v1.PersistentVolumeFilesystem),
			expectedCan: nil,
		},
		{
			name:        "Block volumeMode PV request to BlockProvisoner that returns false",
			provisioner: newTestBlockProvisioner(false),
			claim:       newClaimWithVolumeMode("claim-1", "1-1", "class-1", provisionerName, "", nil, v1.PersistentVolumeBlock),
			expectedCan: fmt.Errorf(blockErrFormat, provisionerName),
		},
		// volumeMode tests for BlockProvisioner that returns true
		{
			name:        "Undefined volumeMode PV request to BlockProvisoner that returns true",
			provisioner: newTestBlockProvisioner(true),
			claim:       newClaim("claim-1", "1-1", "class-1", provisionerName, "", nil),
			expectedCan: nil,
		},
		{
			name:        "FileSystem volumeMode PV request to BlockProvisoner that returns true",
			provisioner: newTestBlockProvisioner(true),
			claim:       newClaimWithVolumeMode("claim-1", "1-1", "class-1", provisionerName, "", nil, v1.PersistentVolumeFilesystem),
			expectedCan: nil,
		},
		{
			name:        "Block volumeMode PV request to BlockProvisioner that returns true",
			provisioner: newTestBlockProvisioner(true),
			claim:       newClaimWithVolumeMode("claim-1", "1-1", "class-1", provisionerName, "", nil, v1.PersistentVolumeBlock),
			expectedCan: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.claim)
			ctrl := newTestProvisionController(client, provisionerName, test.provisioner)

			can := ctrl.canProvision(context.Background(), test.claim)
			if !reflect.DeepEqual(test.expectedCan, can) {
				t.Errorf("expected can provision %v but got %v\n", test.expectedCan, can)
			}
		})
	}
}

func TestControllerSharedInformers(t *testing.T) {
	tests := []struct {
		name            string
		objs            []runtime.Object
		provisionerName string
		expectedVolumes []v1.PersistentVolume
	}{
		{
			name: "provision for claim-1 with v1 storage class",
			objs: []runtime.Object{
				newStorageClassWithReclaimPolicy("class-1", "foo.bar/baz", v1.PersistentVolumeReclaimDelete),
				newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil),
			},
			provisionerName: "foo.bar/baz",
			expectedVolumes: []v1.PersistentVolume{
				*newProvisionedVolumeWithReclaimPolicy(newStorageClassWithReclaimPolicy("class-1", "foo.bar/baz", v1.PersistentVolumeReclaimDelete), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil),
			},
		},
		{
			name: "delete volume-1",
			objs: []runtime.Object{
				newVolume("volume-1", v1.VolumeReleased, v1.PersistentVolumeReclaimDelete, map[string]string{annDynamicallyProvisioned: "foo.bar/baz"}, nil, nil),
			},
			provisionerName: "foo.bar/baz",
			expectedVolumes: []v1.PersistentVolume{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.objs...)

			ctrl, informersFactory := newTestProvisionControllerSharedInformers(client, test.provisionerName,
				newTestProvisioner(), sharedResyncPeriod)
			stopCh := make(chan struct{})
			defer close(stopCh)

			// Run forever...
			go ctrl.Run(context.Background())
			go informersFactory.Start(context.Background().Done())

			// When we shutdown while something is happening the fake client panics
			// with send on closed channel...but the test passed, so ignore
			utilruntime.ReallyCrash = false

			informersFactory.WaitForCacheSync(stopCh)
			time.Sleep(2 * sharedResyncPeriod)

			pvList, _ := client.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
			if (len(test.expectedVolumes) > 0 || len(pvList.Items) > 0) &&
				!reflect.DeepEqual(test.expectedVolumes, pvList.Items) {
				t.Errorf("expected PVs:\n %v\n but got:\n %v\n", test.expectedVolumes, pvList.Items)
			}
		})
	}
}

type testMetrics struct {
	provisioned counts
	deleted     counts
}

type counts map[string]count

type count struct {
	success float64
	failed  float64
}

type testProvisionController struct {
	*ProvisionController
	metrics *metrics.Metrics
}

func (ctrl testProvisionController) getMetrics(t *testing.T) testMetrics {
	var tm testMetrics
	getCounts(t, ctrl.metrics.PersistentVolumeClaimProvisionTotal, &tm.provisioned, true)
	getCounts(t, ctrl.metrics.PersistentVolumeClaimProvisionFailedTotal, &tm.provisioned, false)
	getCounts(t, ctrl.metrics.PersistentVolumeDeleteTotal, &tm.deleted, true)
	getCounts(t, ctrl.metrics.PersistentVolumeDeleteFailedTotal, &tm.deleted, false)
	return tm
}

func getCounts(t *testing.T, vec *prometheus.CounterVec, cts *counts, success bool) {
	metricCh := make(chan prometheus.Metric)
	go func() {
		vec.Collect(metricCh)
		close(metricCh)
	}()
	for metric := range metricCh {
		var m dto.Metric
		err := metric.Write(&m)
		if err != nil {
			t.Fatalf("unexpected error while extracting Prometheus metrics: %v", err)
		}

		// Only initialize the map if we actually have a value.
		if *cts == nil {
			*cts = counts{}
		}

		// We know that our counters have exactly one label.
		count := (*cts)[*m.Label[0].Value]
		if success {
			count.success++
		} else {
			count.failed++
		}
		(*cts)[*m.Label[0].Value] = count
	}
}

func newTestProvisionController(
	client kubernetes.Interface,
	provisionerName string,
	provisioner Provisioner,
	options ...func(*ProvisionController) error,
) testProvisionController {
	m := metrics.New(string(uuid.NewUUID()))
	provisionerOptions := []func(*ProvisionController) error{
		MetricsInstance(m),
		ResyncPeriod(resyncPeriod),
		CreateProvisionedPVInterval(10 * time.Millisecond),
		LeaseDuration(2 * resyncPeriod),
		RenewDeadline(resyncPeriod),
		RetryPeriod(resyncPeriod / 2),
	}
	if len(options) > 0 {
		provisionerOptions = append(provisionerOptions, options...)
	}
	ctrl := NewProvisionController(
		client,
		provisionerName,
		provisioner,
		provisionerOptions...)
	return testProvisionController{
		ProvisionController: ctrl,
		metrics:             &m,
	}
}

func newTestProvisionControllerWithAdditionalNames(
	client kubernetes.Interface,
	provisionerName string,
	provisioner Provisioner,
	additionalProvisionerNames []string,
	options ...func(*ProvisionController) error,
) testProvisionController {
	m := metrics.New(string(uuid.NewUUID()))
	provisionerOptions := []func(*ProvisionController) error{
		MetricsInstance(m),
		ResyncPeriod(resyncPeriod),
		CreateProvisionedPVInterval(10 * time.Millisecond),
		LeaseDuration(2 * resyncPeriod),
		RenewDeadline(resyncPeriod),
		RetryPeriod(resyncPeriod / 2),
		AdditionalProvisionerNames(additionalProvisionerNames),
	}
	if len(options) > 0 {
		provisionerOptions = append(provisionerOptions, options...)
	}
	ctrl := NewProvisionController(
		client,
		provisionerName,
		provisioner,
		provisionerOptions...)
	return testProvisionController{
		ProvisionController: ctrl,
		metrics:             &m,
	}
}

func newTestProvisionControllerSharedInformers(
	client kubernetes.Interface,
	provisionerName string,
	provisioner Provisioner,
	resyncPeriod time.Duration,
) (*ProvisionController, informers.SharedInformerFactory) {

	informerFactory := informers.NewSharedInformerFactory(client, resyncPeriod)
	claimInformer := informerFactory.Core().V1().PersistentVolumeClaims().Informer()
	volumeInformer := informerFactory.Core().V1().PersistentVolumes().Informer()
	classInformer := func() cache.SharedIndexInformer {
		return informerFactory.Storage().V1().StorageClasses().Informer()
	}()

	ctrl := NewProvisionController(
		client,
		provisionerName,
		provisioner,
		ResyncPeriod(resyncPeriod),
		CreateProvisionedPVInterval(10*time.Millisecond),
		LeaseDuration(2*resyncPeriod),
		RenewDeadline(resyncPeriod),
		RetryPeriod(resyncPeriod/2),
		ClaimsInformer(claimInformer),
		VolumesInformer(volumeInformer),
		ClassesInformer(classInformer))

	return ctrl, informerFactory
}

func newStorageClass(name, provisioner string) *storage.StorageClass {
	defaultReclaimPolicy := v1.PersistentVolumeReclaimDelete

	return &storage.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner:   provisioner,
		ReclaimPolicy: &defaultReclaimPolicy,
	}
}

// newStorageClassWithVolumeBindingMode returns the storage class object.
func newStorageClassWithVolumeBindingMode(name, provisioner string, mode *storage.VolumeBindingMode) *storage.StorageClass {
	defaultReclaimPolicy := v1.PersistentVolumeReclaimDelete

	return &storage.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner:       provisioner,
		ReclaimPolicy:     &defaultReclaimPolicy,
		VolumeBindingMode: mode,
	}
}

// newStorageClassWithReclaimPolicy returns the storage class object.
// For Kubernetes version since v1.6.0, it will use the v1 storage class object.
// Once we have tests for v1.6.0, we can add a new function for v1.8.0 newStorageClass since reclaim policy can only be specified since v1.8.0.
func newStorageClassWithReclaimPolicy(name, provisioner string, reclaimPolicy v1.PersistentVolumeReclaimPolicy) *storage.StorageClass {
	return &storage.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner:   provisioner,
		ReclaimPolicy: &reclaimPolicy,
	}
}

func newStorageClassWithAllowedTopologies(name, provisioner string, allowedTopologies []v1.TopologySelectorTerm) *storage.StorageClass {
	defaultReclaimPolicy := v1.PersistentVolumeReclaimDelete

	return &storage.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner:       provisioner,
		ReclaimPolicy:     &defaultReclaimPolicy,
		AllowedTopologies: allowedTopologies,
	}
}

func newClaim(name, claimUID, class, provisioner, volumeName string, annotations map[string]string) *v1.PersistentVolumeClaim {
	claim := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       v1.NamespaceDefault,
			UID:             types.UID(claimUID),
			ResourceVersion: "0",
			Annotations:     map[string]string{},
			SelfLink:        "/api/v1/namespaces/" + v1.NamespaceDefault + "/persistentvolumeclaims/" + name,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce, v1.ReadOnlyMany},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse("1Mi"),
				},
			},
			VolumeName:       volumeName,
			StorageClassName: &class,
		},
		Status: v1.PersistentVolumeClaimStatus{
			Phase: v1.ClaimPending,
		},
	}
	if provisioner != "" {
		claim.Annotations[annBetaStorageProvisioner] = provisioner
	}
	// Allow overwriting of above annotations
	for k, v := range annotations {
		claim.Annotations[k] = v
	}
	return claim
}

func newClaimWithVolumeMode(name, claimUID, class, provisioner, volumeName string, annotations map[string]string, volumeMode v1.PersistentVolumeMode) *v1.PersistentVolumeClaim {
	claim := newClaim(name, claimUID, class, provisioner, volumeName, annotations)
	claim.Spec.VolumeMode = &volumeMode
	return claim
}

func newVolume(name string, phase v1.PersistentVolumePhase, policy v1.PersistentVolumeReclaimPolicy,
	annotations map[string]string, finalizers []string, deletionTimestamp *metav1.Time) *v1.PersistentVolume {
	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Annotations:       annotations,
			Finalizers:        finalizers,
			DeletionTimestamp: deletionTimestamp,
			SelfLink:          "/api/v1/persistentvolumes/" + name,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: policy,
			AccessModes:                   []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce, v1.ReadOnlyMany},
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): resource.MustParse("1Mi"),
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: &v1.NFSVolumeSource{
					Server:   "foo",
					Path:     "bar",
					ReadOnly: false,
				},
			},
		},
		Status: v1.PersistentVolumeStatus{
			Phase: phase,
		},
	}
	return pv
}

func newCSIVolume(name string, phase v1.PersistentVolumePhase, policy v1.PersistentVolumeReclaimPolicy,
	annotations map[string]string, finalizers []string, deletionTimestamp *metav1.Time, provisionerName string) *v1.PersistentVolume {
	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Annotations:       annotations,
			Finalizers:        finalizers,
			DeletionTimestamp: deletionTimestamp,
			SelfLink:          "/api/v1/persistentvolumes/" + name,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: policy,
			AccessModes:                   []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce, v1.ReadOnlyMany},
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): resource.MustParse("1Mi"),
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				CSI: &v1.CSIPersistentVolumeSource{
					Driver:       provisionerName,
					VolumeHandle: "test-volume-id-1",
					ReadOnly:     false,
					FSType:       "ext4",
				},
			},
		},
		Status: v1.PersistentVolumeStatus{
			Phase: phase,
		},
	}
	return pv
}

// newProvisionedVolume returns the volume the test controller should provision
// for the given claim with the given class.
func newProvisionedVolume(storageClass *storage.StorageClass, claim *v1.PersistentVolumeClaim, pvFinalizers []string) *v1.PersistentVolume {
	volume := constructProvisionedVolumeWithoutStorageClassInfo(claim, v1.PersistentVolumeReclaimDelete)

	// pv.Annotations["pv.kubernetes.io/provisioned-by"] MUST be set to name of the external provisioner. This provisioner will be used to delete the volume.
	volume.Annotations = map[string]string{annDynamicallyProvisioned: storageClass.Provisioner}
	// pv.Spec.StorageClassName must be set to the name of the storage class requested by the claim
	volume.Spec.StorageClassName = storageClass.Name
	volume.ObjectMeta.Finalizers = pvFinalizers
	return volume
}

func newProvisionedVolumeWithReclaimPolicy(storageClass *storage.StorageClass, claim *v1.PersistentVolumeClaim, pvFinalizers []string) *v1.PersistentVolume {
	volume := constructProvisionedVolumeWithoutStorageClassInfo(claim, *storageClass.ReclaimPolicy)

	// pv.Annotations["pv.kubernetes.io/provisioned-by"] MUST be set to name of the external provisioner. This provisioner will be used to delete the volume.
	volume.Annotations = map[string]string{annDynamicallyProvisioned: storageClass.Provisioner}
	// pv.Spec.StorageClassName must be set to the name of the storage class requested by the claim
	volume.Spec.StorageClassName = storageClass.Name
	volume.ObjectMeta.Finalizers = pvFinalizers
	return volume
}

func constructProvisionedVolumeWithoutStorageClassInfo(claim *v1.PersistentVolumeClaim, reclaimPolicy v1.PersistentVolumeReclaimPolicy) *v1.PersistentVolume {
	// pv.Spec MUST be set to match requirements in claim.Spec, especially access mode and PV size. The provisioned volume size MUST NOT be smaller than size requested in the claim, however it MAY be larger.
	options := ProvisionOptions{
		StorageClass: &storage.StorageClass{
			ReclaimPolicy: &reclaimPolicy,
		},
		PVName: "pvc-" + string(claim.ObjectMeta.UID),
		PVC:    claim,
	}
	volume, _, _ := newTestProvisioner().Provision(context.Background(), options)

	// pv.Spec.ClaimRef MUST point to the claim that led to its creation (including the claim UID).
	v1.AddToScheme(scheme.Scheme)
	volume.Spec.ClaimRef, _ = ref.GetReference(scheme.Scheme, claim)

	// TODO implement options.ProvisionerSelector parsing
	// pv.Labels MUST be set to match claim.spec.selector. The provisioner MAY add additional labels.

	// TODO addFinalizer is false by default
	// volume.ObjectMeta.Finalizers = append(volume.ObjectMeta.Finalizers, finalizerPV)

	return volume
}

func newNode(nodeName string) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName,
		},
	}
}

type provisionParams struct {
	selectedNode      *v1.Node
	allowedTopologies []v1.TopologySelectorTerm
}

func newTestProvisioner() *testProvisioner {
	return &testProvisioner{make(chan provisionParams, 16)}
}

type testProvisioner struct {
	provisionCalls chan provisionParams
}

var _ Provisioner = &testProvisioner{}

func newTestQualifiedProvisioner(answer bool) *testQualifiedProvisioner {
	return &testQualifiedProvisioner{newTestProvisioner(), answer}
}

type testQualifiedProvisioner struct {
	*testProvisioner
	answer bool
}

var _ Provisioner = &testQualifiedProvisioner{}
var _ Qualifier = &testQualifiedProvisioner{}

func (p *testQualifiedProvisioner) ShouldProvision(ctx context.Context, claim *v1.PersistentVolumeClaim) bool {
	return p.answer
}

func newTestBlockProvisioner(answer bool) *testBlockProvisioner {
	return &testBlockProvisioner{newTestProvisioner(), answer}
}

type testBlockProvisioner struct {
	*testProvisioner
	answer bool
}

var _ Provisioner = &testBlockProvisioner{}
var _ BlockProvisioner = &testBlockProvisioner{}

func (p *testBlockProvisioner) SupportsBlock(ctx context.Context) bool {
	return p.answer
}

func (p *testProvisioner) Provision(ctx context.Context, options ProvisionOptions) (*v1.PersistentVolume, ProvisioningState, error) {
	p.provisionCalls <- provisionParams{
		selectedNode:      options.SelectedNode,
		allowedTopologies: options.StorageClass.AllowedTopologies,
	}

	// Sleep to simulate work done by Provision...for long enough that
	// TestMultipleControllers will consistently fail with lock disabled. If
	// Provision happens too fast, the first controller creates the PV too soon
	// and the next controllers won't call Provision even though they're clearly
	// racing when there's no lock
	time.Sleep(50 * time.Millisecond)

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: &v1.NFSVolumeSource{
					Server:   "foo",
					Path:     "bar",
					ReadOnly: false,
				},
			},
		},
	}

	return pv, ProvisioningFinished, nil
}

func (p *testProvisioner) Delete(ctx context.Context, volume *v1.PersistentVolume) error {
	return nil
}

func newBadTestProvisioner() Provisioner {
	return &badTestProvisioner{}
}

type badTestProvisioner struct {
}

var _ Provisioner = &badTestProvisioner{}

func (p *badTestProvisioner) Provision(ctx context.Context, options ProvisionOptions) (*v1.PersistentVolume, ProvisioningState, error) {
	return nil, ProvisioningFinished, errors.New("fake final error")
}

func (p *badTestProvisioner) Delete(ctx context.Context, volume *v1.PersistentVolume) error {
	return errors.New("fake error")
}

func newTemporaryTestProvisioner() Provisioner {
	return &temporaryTestProvisioner{}
}

type temporaryTestProvisioner struct {
	badTestProvisioner
}

var _ Provisioner = &temporaryTestProvisioner{}

func (p *temporaryTestProvisioner) Provision(ctx context.Context, options ProvisionOptions) (*v1.PersistentVolume, ProvisioningState, error) {
	return nil, ProvisioningInBackground, errors.New("fake error, in progress")
}

func newRescheduleTestProvisioner() Provisioner {
	return &rescheduleTestProvisioner{}
}

type rescheduleTestProvisioner struct {
	badTestProvisioner
}

var _ Provisioner = &rescheduleTestProvisioner{}

func (p *rescheduleTestProvisioner) Provision(ctx context.Context, options ProvisionOptions) (*v1.PersistentVolume, ProvisioningState, error) {
	return nil, ProvisioningReschedule, errors.New("fake error, reschedule")
}

func newNoChangeTestProvisioner() Provisioner {
	return &noChangeTestProvisioner{}
}

type noChangeTestProvisioner struct {
	badTestProvisioner
}

var _ Provisioner = &noChangeTestProvisioner{}

func (p *noChangeTestProvisioner) Provision(ctx context.Context, options ProvisionOptions) (*v1.PersistentVolume, ProvisioningState, error) {
	return nil, ProvisioningNoChange, errors.New("fake error, no change")
}

func newIgnoredProvisioner() Provisioner {
	return &ignoredProvisioner{}
}

type ignoredProvisioner struct {
}

var _ Provisioner = &ignoredProvisioner{}

func (i *ignoredProvisioner) Provision(ctx context.Context, options ProvisionOptions) (*v1.PersistentVolume, ProvisioningState, error) {
	if options.PVC.Name == "claim-2" {
		return nil, ProvisioningFinished, &IgnoredError{"Ignored"}
	}

	return newProvisionedVolume(newStorageClass("class-1", "foo.bar/baz"), newClaim("claim-1", "uid-1-1", "class-1", "foo.bar/baz", "", nil), nil), ProvisioningFinished, nil
}

func (i *ignoredProvisioner) Delete(ctx context.Context, volume *v1.PersistentVolume) error {
	return nil
}

func newProvisioner(t *testing.T, pvName string, returnStatus ProvisioningState, returnError error) Provisioner {
	return &provisioner{
		t:            t,
		pvName:       pvName,
		returnError:  returnError,
		returnStatus: returnStatus,
	}
}

type provisioner struct {
	t            *testing.T
	pvName       string
	returnError  error
	returnStatus ProvisioningState
}

var _ Provisioner = &provisioner{}

func (m *provisioner) Delete(ctx context.Context, pv *v1.PersistentVolume) error {
	return fmt.Errorf("Not implemented")

}

func (m *provisioner) Provision(ctx context.Context, options ProvisionOptions) (*v1.PersistentVolume, ProvisioningState, error) {
	if m.pvName != options.PVName {
		m.t.Errorf("Invalid psrovision call, expected name %q, got %q", m.pvName, options.PVName)
		return nil, ProvisioningFinished, fmt.Errorf("Invalid provision call, expected name %q, got %q", m.pvName, options.PVName)
	}
	klog.Infof("Provision() call")

	if m.returnError == nil {
		pv := &v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name: options.PVName,
			},
			Spec: v1.PersistentVolumeSpec{
				PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
				AccessModes:                   options.PVC.Spec.AccessModes,
				Capacity: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
				},
				PersistentVolumeSource: v1.PersistentVolumeSource{
					NFS: &v1.NFSVolumeSource{
						Server:   "foo",
						Path:     "bar",
						ReadOnly: false,
					},
				},
			},
		}
		return pv, ProvisioningFinished, nil

	}
	return nil, m.returnStatus, m.returnError
}
