# Release notes for v11.1.0

# Changelog since v11.0.0

## Changes by Kind

### Feature

- Added the ability to retry volume provisions that return InvalidArgument at a slower pace than other provisioning failures. ([#186](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/186), [@mdzraf](https://github.com/mdzraf))

### Other (Cleanup or Flake)

- Update module name to sigs.k8s.io/sig-storage-lib-external-provisioner/v11 ([#182](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/182), [@AndrewSirenko](https://github.com/AndrewSirenko))

### Uncategorized

- Update to Kubernetes 1.33 ad go 1.24. ([#188](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/188), [@mdzraf](https://github.com/mdzraf))

## Dependencies

### Added
- cel.dev/expr: v0.16.2
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.24.2](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/tree/detectors/gcp/v1.24.2)
- github.com/alecthomas/kingpin/v2: [v2.4.0](https://github.com/alecthomas/kingpin/tree/v2.4.0)
- github.com/blang/semver/v4: [v4.0.0](https://github.com/blang/semver/tree/v4.0.0)
- github.com/census-instrumentation/opencensus-proto: [v0.4.1](https://github.com/census-instrumentation/opencensus-proto/tree/v0.4.1)
- github.com/cncf/xds/go: [b4127c9](https://github.com/cncf/xds/tree/b4127c9)
- github.com/container-storage-interface/spec: [v1.11.0](https://github.com/container-storage-interface/spec/tree/v1.11.0)
- github.com/envoyproxy/go-control-plane: [v0.13.1](https://github.com/envoyproxy/go-control-plane/tree/v0.13.1)
- github.com/envoyproxy/protoc-gen-validate: [v1.1.0](https://github.com/envoyproxy/protoc-gen-validate/tree/v1.1.0)
- github.com/go-logr/stdr: [v1.2.2](https://github.com/go-logr/stdr/tree/v1.2.2)
- github.com/go-task/slim-sprig/v3: [v3.0.0](https://github.com/go-task/slim-sprig/tree/v3.0.0)
- github.com/golang/glog: [v1.2.2](https://github.com/golang/glog/tree/v1.2.2)
- github.com/jpillora/backoff: [v1.0.0](https://github.com/jpillora/backoff/tree/v1.0.0)
- github.com/klauspost/compress: [v1.18.0](https://github.com/klauspost/compress/tree/v1.18.0)
- github.com/kubernetes-csi/csi-lib-utils: [v0.22.0](https://github.com/kubernetes-csi/csi-lib-utils/tree/v0.22.0)
- github.com/kylelemons/godebug: [v1.1.0](https://github.com/kylelemons/godebug/tree/v1.1.0)
- github.com/planetscale/vtprotobuf: [0393e58](https://github.com/planetscale/vtprotobuf/tree/0393e58)
- github.com/xhit/go-str2duration/v2: [v2.1.0](https://github.com/xhit/go-str2duration/tree/v2.1.0)
- go.opentelemetry.io/auto/sdk: v1.1.0
- go.opentelemetry.io/contrib/detectors/gcp: v1.31.0
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.58.0
- go.opentelemetry.io/otel/metric: v1.33.0
- go.opentelemetry.io/otel/sdk/metric: v1.31.0
- go.opentelemetry.io/otel/sdk: v1.31.0
- go.opentelemetry.io/otel/trace: v1.33.0
- go.opentelemetry.io/otel: v1.33.0
- go.uber.org/automaxprocs: v1.6.0
- go.uber.org/goleak: v1.3.0
- golang.org/x/telemetry: bda5523
- google.golang.org/genproto/googleapis/api: 796eee8
- google.golang.org/genproto/googleapis/rpc: 9240e9c
- google.golang.org/grpc: v1.69.0
- gopkg.in/evanphx/json-patch.v4: v4.12.0
- k8s.io/component-base: v0.33.1
- sigs.k8s.io/randfill: v1.0.0

### Changed
- cloud.google.com/go/compute/metadata: v0.2.3 → v0.5.2
- github.com/NYTimes/gziphandler: [56545f4 → v1.1.1](https://github.com/NYTimes/gziphandler/compare/56545f4...v1.1.1)
- github.com/alecthomas/units: [c3de453 → b94a6e3](https://github.com/alecthomas/units/compare/c3de453...b94a6e3)
- github.com/cespare/xxhash/v2: [v2.1.1 → v2.3.0](https://github.com/cespare/xxhash/compare/v2.1.1...v2.3.0)
- github.com/davecgh/go-spew: [v1.1.1 → d8f796a](https://github.com/davecgh/go-spew/compare/v1.1.1...d8f796a)
- github.com/emicklei/go-restful/v3: [v3.11.0 → v3.12.2](https://github.com/emicklei/go-restful/compare/v3.11.0...v3.12.2)
- github.com/fxamacker/cbor/v2: [v2.6.0 → v2.8.0](https://github.com/fxamacker/cbor/compare/v2.6.0...v2.8.0)
- github.com/go-logr/logr: [v1.4.1 → v1.4.3](https://github.com/go-logr/logr/compare/v1.4.1...v1.4.3)
- github.com/go-openapi/jsonpointer: [v0.19.6 → v0.21.1](https://github.com/go-openapi/jsonpointer/compare/v0.19.6...v0.21.1)
- github.com/go-openapi/jsonreference: [v0.20.2 → v0.21.0](https://github.com/go-openapi/jsonreference/compare/v0.20.2...v0.21.0)
- github.com/go-openapi/swag: [v0.22.3 → v0.23.1](https://github.com/go-openapi/swag/compare/v0.22.3...v0.23.1)
- github.com/google/btree: [v1.0.1 → v1.1.3](https://github.com/google/btree/compare/v1.0.1...v1.1.3)
- github.com/google/gnostic-models: [v0.6.8 → v0.6.9](https://github.com/google/gnostic-models/compare/v0.6.8...v0.6.9)
- github.com/google/go-cmp: [v0.6.0 → v0.7.0](https://github.com/google/go-cmp/compare/v0.6.0...v0.7.0)
- github.com/google/gofuzz: [v1.2.0 → v1.0.0](https://github.com/google/gofuzz/compare/v1.2.0...v1.0.0)
- github.com/google/pprof: [4bb14d4 → d1b30fe](https://github.com/google/pprof/compare/4bb14d4...d1b30fe)
- github.com/google/uuid: [v1.3.0 → v1.6.0](https://github.com/google/uuid/compare/v1.3.0...v1.6.0)
- github.com/gorilla/websocket: [v1.5.0 → e064f32](https://github.com/gorilla/websocket/compare/v1.5.0...e064f32)
- github.com/gregjones/httpcache: [9cad4c3 → 901d907](https://github.com/gregjones/httpcache/compare/9cad4c3...901d907)
- github.com/julienschmidt/httprouter: [v1.2.0 → v1.3.0](https://github.com/julienschmidt/httprouter/compare/v1.2.0...v1.3.0)
- github.com/mailru/easyjson: [v0.7.7 → v0.9.0](https://github.com/mailru/easyjson/compare/v0.7.7...v0.9.0)
- github.com/miekg/dns: [v1.1.29 → v1.1.66](https://github.com/miekg/dns/compare/v1.1.29...v1.1.66)
- github.com/moby/spdystream: [v0.2.0 → v0.5.0](https://github.com/moby/spdystream/compare/v0.2.0...v0.5.0)
- github.com/mwitkow/go-conntrack: [cc309e4 → 2f06839](https://github.com/mwitkow/go-conntrack/compare/cc309e4...2f06839)
- github.com/onsi/ginkgo/v2: [v2.15.0 → v2.21.0](https://github.com/onsi/ginkgo/compare/v2.15.0...v2.21.0)
- github.com/onsi/gomega: [v1.31.0 → v1.35.1](https://github.com/onsi/gomega/compare/v1.31.0...v1.35.1)
- github.com/pmezard/go-difflib: [v1.0.0 → 5d4384e](https://github.com/pmezard/go-difflib/compare/v1.0.0...5d4384e)
- github.com/prometheus/client_golang: [v1.5.1 → v1.22.0](https://github.com/prometheus/client_golang/compare/v1.5.1...v1.22.0)
- github.com/prometheus/client_model: [v0.2.0 → v0.6.2](https://github.com/prometheus/client_model/compare/v0.2.0...v0.6.2)
- github.com/prometheus/common: [v0.9.1 → v0.64.0](https://github.com/prometheus/common/compare/v0.9.1...v0.64.0)
- github.com/prometheus/procfs: [v0.0.8 → v0.16.1](https://github.com/prometheus/procfs/compare/v0.0.8...v0.16.1)
- github.com/rogpeppe/go-internal: [v1.10.0 → v1.13.1](https://github.com/rogpeppe/go-internal/compare/v1.10.0...v1.13.1)
- github.com/stretchr/objx: [v0.5.0 → v0.5.2](https://github.com/stretchr/objx/compare/v0.5.0...v0.5.2)
- github.com/stretchr/testify: [v1.8.4 → v1.10.0](https://github.com/stretchr/testify/compare/v1.8.4...v1.10.0)
- github.com/yuin/goldmark: [v1.2.1 → v1.4.13](https://github.com/yuin/goldmark/compare/v1.2.1...v1.4.13)
- golang.org/x/crypto: v0.21.0 → v0.38.0
- golang.org/x/mod: v0.15.0 → v0.24.0
- golang.org/x/net: v0.23.0 → v0.40.0
- golang.org/x/oauth2: v0.10.0 → v0.30.0
- golang.org/x/sync: 886fb93 → v0.14.0
- golang.org/x/sys: v0.18.0 → v0.33.0
- golang.org/x/term: v0.18.0 → v0.32.0
- golang.org/x/text: v0.14.0 → v0.25.0
- golang.org/x/time: v0.3.0 → v0.11.0
- golang.org/x/tools: v0.18.0 → v0.33.0
- golang.org/x/xerrors: 04be3eb → 5ec99f8
- google.golang.org/protobuf: v1.33.0 → v1.36.6
- k8s.io/api: v0.30.0 → v0.33.1
- k8s.io/apimachinery: v0.30.0 → v0.33.1
- k8s.io/client-go: v0.30.0 → v0.33.1
- k8s.io/gengo/v2: 51d4e06 → a7b603a
- k8s.io/klog/v2: v2.120.1 → v2.130.1
- k8s.io/kube-openapi: 70dd376 → c8a335a
- k8s.io/utils: 3b25d92 → 0f33e8f
- sigs.k8s.io/json: bc3834c → cfa47c3
- sigs.k8s.io/structured-merge-diff/v4: v4.4.1 → v4.7.0
- sigs.k8s.io/yaml: v1.3.0 → v1.4.0

### Removed
- cloud.google.com/go/compute: v1.20.1
- github.com/alecthomas/template: [fb15b89](https://github.com/alecthomas/template/tree/fb15b89)
- github.com/asaskevich/govalidator: [f61b66f](https://github.com/asaskevich/govalidator/tree/f61b66f)
- github.com/creack/pty: [v1.1.9](https://github.com/creack/pty/tree/v1.1.9)
- github.com/evanphx/json-patch: [v5.6.0+incompatible](https://github.com/evanphx/json-patch/tree/v5.6.0)
- github.com/go-kit/kit: [v0.9.0](https://github.com/go-kit/kit/tree/v0.9.0)
- github.com/go-logfmt/logfmt: [v0.4.0](https://github.com/go-logfmt/logfmt/tree/v0.4.0)
- github.com/go-stack/stack: [v1.8.0](https://github.com/go-stack/stack/tree/v1.8.0)
- github.com/go-task/slim-sprig: [52ccab3](https://github.com/go-task/slim-sprig/tree/52ccab3)
- github.com/golang/groupcache: [41bb18b](https://github.com/golang/groupcache/tree/41bb18b)
- github.com/imdario/mergo: [v0.3.6](https://github.com/imdario/mergo/tree/v0.3.6)
- github.com/konsorten/go-windows-terminal-sequences: [v1.0.1](https://github.com/konsorten/go-windows-terminal-sequences/tree/v1.0.1)
- github.com/kr/logfmt: [b84e30a](https://github.com/kr/logfmt/tree/b84e30a)
- github.com/kr/pty: [v1.1.1](https://github.com/kr/pty/tree/v1.1.1)
- github.com/matttproud/golang_protobuf_extensions: [v1.0.1](https://github.com/matttproud/golang_protobuf_extensions/tree/v1.0.1)
- github.com/sirupsen/logrus: [v1.4.2](https://github.com/sirupsen/logrus/tree/v1.4.2)
- google.golang.org/appengine: v1.6.7
- gopkg.in/alecthomas/kingpin.v2: v2.2.6
