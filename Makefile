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

GOBIN=$(shell pwd)/hack/tools/bin

all: install-tools verify golangci-lint test

dep:
	go mod tidy
	cd hack/tools && go mod tidy

install-tools: $(GOBIN)
	cd hack/tools \
		&& GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint \
		&& go build -o "${GOBIN}/logcheck.so" -buildmode=plugin sigs.k8s.io/logtools/logcheck/plugin

$(GOBIN):
	mkdir $@

golangci-lint:
	LOGCHECK_CONFIG="hack/logcheck.conf" "${GOBIN}/golangci-lint" run

verify: dep
	PATH=$$(go env GOPATH)/bin:$$PATH repo-infra/verify/verify-go-src.sh -v
	PATH=$$(go env GOPATH)/bin:$$PATH repo-infra/verify/verify-boilerplate.sh

test: dep
	go test ./controller -v
	go test ./allocator -v

clean:
	rm -rf ./test/e2e/kubernetes
