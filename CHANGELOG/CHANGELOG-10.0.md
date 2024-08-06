# Release notes for v10.0.1

## Changes by Kind

### Bug or Regression

- Fixed removal of PV protection finalizer. PVs are no longer Terminating forever after PVC deletion. ([#174](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/174), [@jsafrane](https://github.com/jsafrane))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v10.0.0

# Changelog since v9.0.3

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)
- Using patch to update finalizers. Any external-provisioner now needs permission rules to patch PersistentVolumes. Please update RBAC rules of your provisioner. ([#164](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/164), [@carlory](https://github.com/carlory))

## Changes by Kind

### Feature

- Contextual logging added. Some function arguments have been modified. ([#154](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/154), [@bells17](https://github.com/bells17))

### Uncategorized

- Added debug logs to shouldDelete function ([#146](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/146), [@amacaskill](https://github.com/amacaskill))
- Update to Kubernetes 1.30 and go 1.22. Removed gometalinter. ([#167](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/167), [@jsafrane](https://github.com/jsafrane))

## Dependencies

### Added
- cloud.google.com/go/compute/metadata: v0.2.3
- cloud.google.com/go/compute: v1.20.1
- github.com/fxamacker/cbor/v2: [v2.6.0](https://github.com/fxamacker/cbor/v2/tree/v2.6.0)
- github.com/go-task/slim-sprig: [52ccab3](https://github.com/go-task/slim-sprig/tree/52ccab3)
- github.com/google/gnostic-models: [v0.6.8](https://github.com/google/gnostic-models/tree/v0.6.8)
- github.com/gorilla/websocket: [v1.5.0](https://github.com/gorilla/websocket/tree/v1.5.0)
- github.com/x448/float16: [v0.8.4](https://github.com/x448/float16/tree/v0.8.4)
- k8s.io/gengo/v2: 51d4e06

### Changed
- github.com/emicklei/go-restful/v3: [v3.9.0 → v3.11.0](https://github.com/emicklei/go-restful/v3/compare/v3.9.0...v3.11.0)
- github.com/evanphx/json-patch: [v4.12.0+incompatible → v5.6.0+incompatible](https://github.com/evanphx/json-patch/compare/v4.12.0...v5.6.0)
- github.com/go-logr/logr: [v1.2.3 → v1.4.1](https://github.com/go-logr/logr/compare/v1.2.3...v1.4.1)
- github.com/go-openapi/jsonpointer: [v0.19.5 → v0.19.6](https://github.com/go-openapi/jsonpointer/compare/v0.19.5...v0.19.6)
- github.com/go-openapi/jsonreference: [v0.20.0 → v0.20.2](https://github.com/go-openapi/jsonreference/compare/v0.20.0...v0.20.2)
- github.com/go-openapi/swag: [v0.19.14 → v0.22.3](https://github.com/go-openapi/swag/compare/v0.19.14...v0.22.3)
- github.com/golang/protobuf: [v1.5.2 → v1.5.4](https://github.com/golang/protobuf/compare/v1.5.2...v1.5.4)
- github.com/google/go-cmp: [v0.5.9 → v0.6.0](https://github.com/google/go-cmp/compare/v0.5.9...v0.6.0)
- github.com/google/gofuzz: [v1.1.0 → v1.2.0](https://github.com/google/gofuzz/compare/v1.1.0...v1.2.0)
- github.com/google/pprof: [1a94d86 → 4bb14d4](https://github.com/google/pprof/compare/1a94d86...4bb14d4)
- github.com/google/uuid: [v1.1.2 → v1.3.0](https://github.com/google/uuid/compare/v1.1.2...v1.3.0)
- github.com/kr/pretty: [v0.2.0 → v0.3.1](https://github.com/kr/pretty/compare/v0.2.0...v0.3.1)
- github.com/mailru/easyjson: [v0.7.6 → v0.7.7](https://github.com/mailru/easyjson/compare/v0.7.6...v0.7.7)
- github.com/onsi/ginkgo/v2: [v2.4.0 → v2.15.0](https://github.com/onsi/ginkgo/v2/compare/v2.4.0...v2.15.0)
- github.com/onsi/gomega: [v1.23.0 → v1.31.0](https://github.com/onsi/gomega/compare/v1.23.0...v1.31.0)
- github.com/rogpeppe/go-internal: [v1.3.0 → v1.10.0](https://github.com/rogpeppe/go-internal/compare/v1.3.0...v1.10.0)
- github.com/stretchr/objx: [v0.1.1 → v0.5.0](https://github.com/stretchr/objx/compare/v0.1.1...v0.5.0)
- github.com/stretchr/testify: [v1.8.0 → v1.8.4](https://github.com/stretchr/testify/compare/v1.8.0...v1.8.4)
- golang.org/x/crypto: v0.1.0 → v0.21.0
- golang.org/x/mod: 86c51ed → v0.15.0
- golang.org/x/net: 1e63c2f → v0.23.0
- golang.org/x/oauth2: ee48083 → v0.10.0
- golang.org/x/sys: v0.3.0 → v0.18.0
- golang.org/x/term: v0.3.0 → v0.18.0
- golang.org/x/text: v0.5.0 → v0.14.0
- golang.org/x/time: 90d013b → v0.3.0
- golang.org/x/tools: v0.1.12 → v0.18.0
- golang.org/x/xerrors: 5ec99f8 → 04be3eb
- google.golang.org/protobuf: v1.28.1 → v1.33.0
- gopkg.in/check.v1: 8fa4692 → 10cb982
- k8s.io/api: v0.26.0 → v0.30.0
- k8s.io/apimachinery: v0.26.0 → v0.30.0
- k8s.io/client-go: v0.26.0 → v0.30.0
- k8s.io/klog/v2: v2.80.1 → v2.120.1
- k8s.io/kube-openapi: 172d655 → 70dd376
- k8s.io/utils: 1a15be2 → 3b25d92
- sigs.k8s.io/json: f223a00 → bc3834c
- sigs.k8s.io/structured-merge-diff/v4: v4.2.3 → v4.4.1

### Removed
- cloud.google.com/go/bigquery: v1.8.0
- cloud.google.com/go/datastore: v1.1.0
- cloud.google.com/go/pubsub: v1.3.1
- cloud.google.com/go/storage: v1.10.0
- cloud.google.com/go: v0.65.0
- dmitri.shuralyov.com/gpu/mtl: 666a987
- github.com/BurntSushi/toml: [v0.3.1](https://github.com/BurntSushi/toml/tree/v0.3.1)
- github.com/BurntSushi/xgb: [27f1227](https://github.com/BurntSushi/xgb/tree/27f1227)
- github.com/PuerkitoBio/purell: [v1.1.1](https://github.com/PuerkitoBio/purell/tree/v1.1.1)
- github.com/PuerkitoBio/urlesc: [de5bf2a](https://github.com/PuerkitoBio/urlesc/tree/de5bf2a)
- github.com/census-instrumentation/opencensus-proto: [v0.2.1](https://github.com/census-instrumentation/opencensus-proto/tree/v0.2.1)
- github.com/chzyer/logex: [v1.1.10](https://github.com/chzyer/logex/tree/v1.1.10)
- github.com/chzyer/readline: [2972be2](https://github.com/chzyer/readline/tree/2972be2)
- github.com/chzyer/test: [a1ea475](https://github.com/chzyer/test/tree/a1ea475)
- github.com/client9/misspell: [v0.3.4](https://github.com/client9/misspell/tree/v0.3.4)
- github.com/cncf/udpa/go: [269d4d4](https://github.com/cncf/udpa/go/tree/269d4d4)
- github.com/docopt/docopt-go: [ee0de3b](https://github.com/docopt/docopt-go/tree/ee0de3b)
- github.com/elazarl/goproxy: [947c36d](https://github.com/elazarl/goproxy/tree/947c36d)
- github.com/envoyproxy/go-control-plane: [v0.9.4](https://github.com/envoyproxy/go-control-plane/tree/v0.9.4)
- github.com/envoyproxy/protoc-gen-validate: [v0.1.0](https://github.com/envoyproxy/protoc-gen-validate/tree/v0.1.0)
- github.com/go-gl/glfw/v3.3/glfw: [6f7a984](https://github.com/go-gl/glfw/v3.3/glfw/tree/6f7a984)
- github.com/go-gl/glfw: [e6da0ac](https://github.com/go-gl/glfw/tree/e6da0ac)
- github.com/golang/glog: [23def4e](https://github.com/golang/glog/tree/23def4e)
- github.com/golang/mock: [v1.4.4](https://github.com/golang/mock/tree/v1.4.4)
- github.com/google/gnostic: [v0.5.7-v3refs](https://github.com/google/gnostic/tree/v0.5.7-v3refs)
- github.com/google/martian/v3: [v3.0.0](https://github.com/google/martian/v3/tree/v3.0.0)
- github.com/google/martian: [v2.1.0+incompatible](https://github.com/google/martian/tree/v2.1.0)
- github.com/google/renameio: [v0.1.0](https://github.com/google/renameio/tree/v0.1.0)
- github.com/googleapis/gax-go/v2: [v2.0.5](https://github.com/googleapis/gax-go/v2/tree/v2.0.5)
- github.com/hashicorp/golang-lru: [v0.5.1](https://github.com/hashicorp/golang-lru/tree/v0.5.1)
- github.com/ianlancetaylor/demangle: [5e5cf60](https://github.com/ianlancetaylor/demangle/tree/5e5cf60)
- github.com/jstemmer/go-junit-report: [v0.9.1](https://github.com/jstemmer/go-junit-report/tree/v0.9.1)
- github.com/mitchellh/mapstructure: [v1.1.2](https://github.com/mitchellh/mapstructure/tree/v1.1.2)
- github.com/niemeyer/pretty: [a10e7ca](https://github.com/niemeyer/pretty/tree/a10e7ca)
- github.com/stoewer/go-strcase: [v1.2.0](https://github.com/stoewer/go-strcase/tree/v1.2.0)
- go.opencensus.io: v0.22.4
- golang.org/x/exp: 6cc2880
- golang.org/x/image: cff245a
- golang.org/x/lint: 738671d
- golang.org/x/mobile: d2bd2a2
- google.golang.org/api: v0.30.0
- google.golang.org/genproto: 1ed22bb
- google.golang.org/grpc: v1.31.0
- gopkg.in/errgo.v2: v2.1.0
- honnef.co/go/tools: v0.0.1-2020.1.4
- k8s.io/gengo: 485abfe
- rsc.io/binaryregexp: v0.2.0
- rsc.io/quote/v3: v3.1.0
- rsc.io/sampler: v1.3.0
 v0.3.0
- golang.org/x/tools: v0.1.12 → v0.18.0
- golang.org/x/xerrors: 5ec99f8 → 04be3eb
- google.golang.org/protobuf: v1.28.1 → v1.33.0
- gopkg.in/check.v1: 8fa4692 → 10cb982
- k8s.io/api: v0.26.0 → v0.30.0
- k8s.io/apimachinery: v0.26.0 → v0.30.0
- k8s.io/client-go: v0.26.0 → v0.30.0
- k8s.io/klog/v2: v2.80.1 → v2.120.1
- k8s.io/kube-openapi: 172d655 → 70dd376
- k8s.io/utils: 1a15be2 → 3b25d92
- sigs.k8s.io/json: f223a00 → bc3834c
- sigs.k8s.io/structured-merge-diff/v4: v4.2.3 → v4.4.1

### Removed
- cloud.google.com/go/bigquery: v1.8.0
- cloud.google.com/go/datastore: v1.1.0
- cloud.google.com/go/pubsub: v1.3.1
- cloud.google.com/go/storage: v1.10.0
- cloud.google.com/go: v0.65.0
- dmitri.shuralyov.com/gpu/mtl: 666a987
- github.com/BurntSushi/toml: [v0.3.1](https://github.com/BurntSushi/toml/tree/v0.3.1)
- github.com/BurntSushi/xgb: [27f1227](https://github.com/BurntSushi/xgb/tree/27f1227)
- github.com/PuerkitoBio/purell: [v1.1.1](https://github.com/PuerkitoBio/purell/tree/v1.1.1)
- github.com/PuerkitoBio/urlesc: [de5bf2a](https://github.com/PuerkitoBio/urlesc/tree/de5bf2a)
- github.com/census-instrumentation/opencensus-proto: [v0.2.1](https://github.com/census-instrumentation/opencensus-proto/tree/v0.2.1)
- github.com/chzyer/logex: [v1.1.10](https://github.com/chzyer/logex/tree/v1.1.10)
- github.com/chzyer/readline: [2972be2](https://github.com/chzyer/readline/tree/2972be2)
- github.com/chzyer/test: [a1ea475](https://github.com/chzyer/test/tree/a1ea475)
- github.com/client9/misspell: [v0.3.4](https://github.com/client9/misspell/tree/v0.3.4)
- github.com/cncf/udpa/go: [269d4d4](https://github.com/cncf/udpa/go/tree/269d4d4)
- github.com/docopt/docopt-go: [ee0de3b](https://github.com/docopt/docopt-go/tree/ee0de3b)
- github.com/elazarl/goproxy: [947c36d](https://github.com/elazarl/goproxy/tree/947c36d)
- github.com/envoyproxy/go-control-plane: [v0.9.4](https://github.com/envoyproxy/go-control-plane/tree/v0.9.4)
- github.com/envoyproxy/protoc-gen-validate: [v0.1.0](https://github.com/envoyproxy/protoc-gen-validate/tree/v0.1.0)
- github.com/go-gl/glfw/v3.3/glfw: [6f7a984](https://github.com/go-gl/glfw/v3.3/glfw/tree/6f7a984)
- github.com/go-gl/glfw: [e6da0ac](https://github.com/go-gl/glfw/tree/e6da0ac)
- github.com/golang/glog: [23def4e](https://github.com/golang/glog/tree/23def4e)
- github.com/golang/mock: [v1.4.4](https://github.com/golang/mock/tree/v1.4.4)
- github.com/google/gnostic: [v0.5.7-v3refs](https://github.com/google/gnostic/tree/v0.5.7-v3refs)
- github.com/google/martian/v3: [v3.0.0](https://github.com/google/martian/v3/tree/v3.0.0)
- github.com/google/martian: [v2.1.0+incompatible](https://github.com/google/martian/tree/v2.1.0)
- github.com/google/renameio: [v0.1.0](https://github.com/google/renameio/tree/v0.1.0)
- github.com/googleapis/gax-go/v2: [v2.0.5](https://github.com/googleapis/gax-go/v2/tree/v2.0.5)
- github.com/hashicorp/golang-lru: [v0.5.1](https://github.com/hashicorp/golang-lru/tree/v0.5.1)
- github.com/ianlancetaylor/demangle: [5e5cf60](https://github.com/ianlancetaylor/demangle/tree/5e5cf60)
- github.com/jstemmer/go-junit-report: [v0.9.1](https://github.com/jstemmer/go-junit-report/tree/v0.9.1)
- github.com/mitchellh/mapstructure: [v1.1.2](https://github.com/mitchellh/mapstructure/tree/v1.1.2)
- github.com/niemeyer/pretty: [a10e7ca](https://github.com/niemeyer/pretty/tree/a10e7ca)
- github.com/stoewer/go-strcase: [v1.2.0](https://github.com/stoewer/go-strcase/tree/v1.2.0)
- go.opencensus.io: v0.22.4
- golang.org/x/exp: 6cc2880
- golang.org/x/image: cff245a
- golang.org/x/lint: 738671d
- golang.org/x/mobile: d2bd2a2
- google.golang.org/api: v0.30.0
- google.golang.org/genproto: 1ed22bb
- google.golang.org/grpc: v1.31.0
- gopkg.in/errgo.v2: v2.1.0
- honnef.co/go/tools: v0.0.1-2020.1.4
- k8s.io/gengo: 485abfe
- rsc.io/binaryregexp: v0.2.0
- rsc.io/quote/v3: v3.1.0
- rsc.io/sampler: v1.3.0
