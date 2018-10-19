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

all: verify test

dep:
	cp .Gopkg.toml Gopkg.toml
	-dep init
	dep ensure

verify: dep
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	repo-infra/verify/verify-go-src.sh -v
	repo-infra/verify/verify-boilerplate.sh

test: dep
	go test ./controller
	go test ./allocator

clean:
	rm -rf ./vendor
	rm -rf ./Gopkg.toml
	rm -rf ./Gopkg.lock
	rm -rf ./examples/hostpath-provisioner/vendor
	rm -rf ./examples/hostpath-provisioner/Gopkg.lock
	rm -rf ./test/e2e/kubernetes
