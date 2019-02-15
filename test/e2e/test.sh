#!/usr/bin/env bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -x
PS4='+\t '

pushd $GOPATH/src/sigs.k8s.io/sig-storage-lib-external-provisioner/test/e2e
if [ ! -f dind-cluster-v1.13.sh ]; then
  wget https://github.com/kubernetes-sigs/kubeadm-dind-cluster/releases/download/v0.1.0/dind-cluster-v1.13.sh
  chmod +x dind-cluster-v1.13.sh
fi
./dind-cluster-v1.13.sh up

export PATH="$HOME/.kubeadm-dind-cluster:$PATH"

pushd $GOPATH/src/sigs.k8s.io/sig-storage-lib-external-provisioner/examples/hostpath-provisioner
make image
docker save hostpath-provisioner | docker exec -i kube-node-1 docker load
docker save hostpath-provisioner | docker exec -i kube-node-2 docker load
kubectl create -f ./rbac.yaml
kubectl create -f ./pod.yaml
kubectl create -f ./class.yaml
kubectl create -f ./claim.yaml
kubectl create -f ./test-pod.yaml
timeout 30 bash -c "until kubectl get pod test-pod -o=jsonpath='{.status.phase}' | grep -E 'Succeeded|Failed'; do sleep 1; done"
kubectl describe pod test-pod
kubectl describe pod hostpath-provisioner
kubectl logs hostpath-provisioner
kubectl describe pvc
kubectl describe pv
if [ $? == 0 ] && kubectl get pod test-pod -o=jsonpath='{.status.phase}' | grep -q Succeeded; then
  #./dind-cluster-v1.13.sh down
  exit 0
fi
#./dind-cluster-v1.13.sh down
exit 1
