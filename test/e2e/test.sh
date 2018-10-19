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

curl -L https://dl.k8s.io/v1.12.0/kubernetes-server-linux-amd64.tar.gz | tar xz
tar xz -f kubernetes/kubernetes-src.tar.gz -C kubernetes

kubernetes/hack/install-etcd.sh
export PATH=$PWD/kubernetes/third_party/etcd:${PATH}

export KUBECTL=$PWD/kubernetes/server/bin/kubectl

sudo "KUBECTL=$KUBECTL" "PATH=$PATH" kubernetes/hack/local-up-cluster.sh -o kubernetes/server/bin > /tmp/local-up-cluster.log 2>&1 &

timeout 60 grep -q "Local Kubernetes cluster is running." <(tail -f /tmp/local-up-cluster.log)
code=$?
if [ $code != 0 ]; then
  exit 1
fi

pushd $GOPATH/src/github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/examples/hostpath-provisioner
make image
$KUBECTL create -f ./rbac.yaml
$KUBECTL create -f ./pod.yaml
$KUBECTL create -f ./class.yaml
$KUBECTL create -f ./claim.yaml
$KUBECTL create -f ./test-pod.yaml
timeout 10 bash -c "until $KUBECTL get pod test-pod -o=jsonpath='{.status.phase}' | grep -E 'Succeeded|Failed'; do sleep 1; done"
if [ $? == 0 ] && $KUBECTL get pod test-pod -o=jsonpath='{.status.phase}' | grep -q Succeeded; then
  exit 0
fi
exit 1
