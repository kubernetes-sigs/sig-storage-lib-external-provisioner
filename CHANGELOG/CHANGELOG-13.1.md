# Release notes for v13.1.0

# Changelog since v13.0.0

## Changes by Kind

### Feature

- Provisioners can return `VolumeInUseError` from `Delete` to signal that a volume is still in use; the controller then retries without emitting a `VolumeFailedDelete` warning. ([#206](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/206), [@WanzenBug](https://github.com/WanzenBug))
- Updated dependencies to Kubernetes 1.37 ([#204](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/204), [@jsafrane](https://github.com/jsafrane))

## Dependencies

### Added
- cloud.google.com/go/auth: v0.18.2
- github.com/felixge/httpsnoop: [v1.0.4](https://github.com/felixge/httpsnoop/tree/v1.0.4)
- github.com/go-openapi/swag/jsonutils/fixtures_test: [v0.29.1](https://github.com/go-openapi/swag/tree/jsonutils/fixtures_test/v0.29.1)
- github.com/go-openapi/swag/pools: [v0.29.1](https://github.com/go-openapi/swag/tree/pools/v0.29.1)
- github.com/go-openapi/testify/enable/yaml/v2: [v2.6.1](https://github.com/go-openapi/testify/tree/enable/yaml/v2/v2.6.1)
- github.com/go-openapi/testify/v2: [v2.6.1](https://github.com/go-openapi/testify/tree/v2.6.1)
- github.com/golang-jwt/jwt/v5: [v5.3.1](https://github.com/golang-jwt/jwt/tree/v5.3.1)
- github.com/google/s2a-go: [v0.1.9](https://github.com/google/s2a-go/tree/v0.1.9)
- github.com/googleapis/enterprise-certificate-proxy: [v0.3.11](https://github.com/googleapis/enterprise-certificate-proxy/tree/v0.3.11)
- github.com/googleapis/gax-go/v2: [v2.17.0](https://github.com/googleapis/gax-go/tree/v2.17.0)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.61.0
- k8s.io/streaming: v0.37.0

### Changed
- cel.dev/expr: v0.24.0 → v0.25.2
- cloud.google.com/go/compute/metadata: v0.7.0 → v0.9.0
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.29.0 → v1.33.0](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/compare/detectors/gcp/v1.29.0...detectors/gcp/v1.33.0)
- github.com/alecthomas/units: [b94a6e3 → 0f3dac3](https://github.com/alecthomas/units/compare/b94a6e3...0f3dac3)
- github.com/cncf/xds/go: [2ac532f → dba9d58](https://github.com/cncf/xds/compare/2ac532f...dba9d58)
- github.com/container-storage-interface/spec: [v1.11.0 → v1.13.0](https://github.com/container-storage-interface/spec/compare/v1.11.0...v1.13.0)
- github.com/envoyproxy/go-control-plane/envoy: [v1.32.4 → v1.37.0](https://github.com/envoyproxy/go-control-plane/compare/envoy/v1.32.4...envoy/v1.37.0)
- github.com/envoyproxy/go-control-plane: [v0.13.4 → v0.14.0](https://github.com/envoyproxy/go-control-plane/compare/v0.13.4...v0.14.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.2.1 → v1.3.3](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.2.1...v1.3.3)
- github.com/fxamacker/cbor/v2: [v2.9.0 → v2.9.3](https://github.com/fxamacker/cbor/compare/v2.9.0...v2.9.3)
- github.com/go-jose/go-jose/v4: [v4.1.1 → v4.1.4](https://github.com/go-jose/go-jose/compare/v4.1.1...v4.1.4)
- github.com/go-logr/logr: [v1.4.3 → v1.4.4](https://github.com/go-logr/logr/compare/v1.4.3...v1.4.4)
- github.com/go-openapi/jsonpointer: [v0.22.0 → v1.0.0](https://github.com/go-openapi/jsonpointer/compare/v0.22.0...v1.0.0)
- github.com/go-openapi/jsonreference: [v0.21.1 → v1.0.1](https://github.com/go-openapi/jsonreference/compare/v0.21.1...v1.0.1)
- github.com/go-openapi/swag/cmdutils: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/cmdutils/v0.24.0...cmdutils/v0.29.1)
- github.com/go-openapi/swag/conv: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/conv/v0.24.0...conv/v0.29.1)
- github.com/go-openapi/swag/fileutils: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/fileutils/v0.24.0...fileutils/v0.29.1)
- github.com/go-openapi/swag/jsonutils: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/jsonutils/v0.24.0...jsonutils/v0.29.1)
- github.com/go-openapi/swag/loading: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/loading/v0.24.0...loading/v0.29.1)
- github.com/go-openapi/swag/mangling: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/mangling/v0.24.0...mangling/v0.29.1)
- github.com/go-openapi/swag/netutils: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/netutils/v0.24.0...netutils/v0.29.1)
- github.com/go-openapi/swag/stringutils: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/stringutils/v0.24.0...stringutils/v0.29.1)
- github.com/go-openapi/swag/typeutils: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/typeutils/v0.24.0...typeutils/v0.29.1)
- github.com/go-openapi/swag/yamlutils: [v0.24.0 → v0.29.1](https://github.com/go-openapi/swag/compare/yamlutils/v0.24.0...yamlutils/v0.29.1)
- github.com/go-openapi/swag: [v0.24.1 → v0.29.1](https://github.com/go-openapi/swag/compare/v0.24.1...v0.29.1)
- github.com/google/gnostic-models: [v0.7.0 → v0.7.1](https://github.com/google/gnostic-models/compare/v0.7.0...v0.7.1)
- github.com/klauspost/compress: [v1.18.0 → v1.19.1](https://github.com/klauspost/compress/compare/v1.18.0...v1.19.1)
- github.com/kubernetes-csi/csi-lib-utils: [v0.22.0 → v0.25.0](https://github.com/kubernetes-csi/csi-lib-utils/compare/v0.22.0...v0.25.0)
- github.com/miekg/dns: [v1.1.68 → v1.1.73](https://github.com/miekg/dns/compare/v1.1.68...v1.1.73)
- github.com/moby/spdystream: [v0.5.0 → v0.5.1](https://github.com/moby/spdystream/compare/v0.5.0...v0.5.1)
- github.com/prometheus/client_golang: [v1.23.0 → v1.24.1](https://github.com/prometheus/client_golang/compare/v1.23.0...v1.24.1)
- github.com/prometheus/client_model: [v0.6.2 → v0.6.3](https://github.com/prometheus/client_model/compare/v0.6.2...v0.6.3)
- github.com/prometheus/common: [v0.66.0 → v0.71.0](https://github.com/prometheus/common/compare/v0.66.0...v0.71.0)
- github.com/prometheus/procfs: [v0.17.0 → v0.22.0](https://github.com/prometheus/procfs/compare/v0.17.0...v0.22.0)
- github.com/rogpeppe/go-internal: [v1.13.1 → v1.14.1](https://github.com/rogpeppe/go-internal/compare/v1.13.1...v1.14.1)
- github.com/spf13/pflag: [v1.0.6 → v1.0.10](https://github.com/spf13/pflag/compare/v1.0.6...v1.0.10)
- github.com/spiffe/go-spiffe/v2: [v2.5.0 → v2.7.0](https://github.com/spiffe/go-spiffe/compare/v2.5.0...v2.7.0)
- github.com/stretchr/objx: [v0.5.2 → v0.5.3](https://github.com/stretchr/objx/compare/v0.5.2...v0.5.3)
- github.com/stretchr/testify: [v1.11.1 → v1.12.1](https://github.com/stretchr/testify/compare/v1.11.1...v1.12.1)
- go.opentelemetry.io/auto/sdk: v1.1.0 → v1.2.1
- go.opentelemetry.io/contrib/detectors/gcp: v1.36.0 → v1.44.0
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.58.0 → v0.71.0
- go.opentelemetry.io/otel/metric: v1.37.0 → v1.46.0
- go.opentelemetry.io/otel/sdk/metric: v1.37.0 → v1.44.0
- go.opentelemetry.io/otel/sdk: v1.37.0 → v1.44.0
- go.opentelemetry.io/otel/trace: v1.37.0 → v1.46.0
- go.opentelemetry.io/otel: v1.37.0 → v1.46.0
- go.yaml.in/yaml/v2: v2.4.2 → v2.4.4
- go.yaml.in/yaml/v3: v3.0.4 → v3.0.5
- golang.org/x/crypto: v0.41.0 → v0.55.0
- golang.org/x/mod: v0.27.0 → v0.38.0
- golang.org/x/net: v0.43.0 → v0.58.0
- golang.org/x/oauth2: v0.30.0 → v0.36.0
- golang.org/x/sync: v0.16.0 → v0.22.0
- golang.org/x/sys: v0.35.0 → v0.47.0
- golang.org/x/term: v0.34.0 → v0.45.0
- golang.org/x/text: v0.28.0 → v0.41.0
- golang.org/x/time: v0.12.0 → v0.15.0
- golang.org/x/tools: v0.36.0 → v0.48.0
- gonum.org/v1/gonum: v0.16.0 → v0.17.0
- google.golang.org/genproto/googleapis/api: 8d1bb00 → 3dc84a4
- google.golang.org/genproto/googleapis/rpc: ef028d9 → da73d73
- google.golang.org/grpc: v1.75.0 → v1.83.2
- google.golang.org/protobuf: v1.36.8 → v1.36.12
- k8s.io/api: v0.34.0 → v0.37.0
- k8s.io/apimachinery: v0.34.0 → v0.37.0
- k8s.io/client-go: v0.34.0 → v0.37.0
- k8s.io/component-base: v0.33.1 → v0.37.0
- k8s.io/gengo/v2: 85fd79d → ec3ebc5
- k8s.io/klog/v2: v2.130.1 → v2.140.0
- k8s.io/kube-openapi: 7fc2783 → be32def
- k8s.io/utils: 0af2bda → cf1189d
- sigs.k8s.io/structured-merge-diff/v6: v6.3.0 → v6.4.2

### Removed
- github.com/armon/go-socks5: [e753329](https://github.com/armon/go-socks5/tree/e753329)
- github.com/go-openapi/swag/jsonname: [v0.24.0](https://github.com/go-openapi/swag/tree/jsonname/v0.24.0)
- github.com/go-task/slim-sprig/v3: [v3.0.0](https://github.com/go-task/slim-sprig/tree/v3.0.0)
- github.com/gogo/protobuf: [v1.3.2](https://github.com/gogo/protobuf/tree/v1.3.2)
- github.com/google/pprof: [d1b30fe](https://github.com/google/pprof/tree/d1b30fe)
- github.com/grafana/regexp: [a468a5b](https://github.com/grafana/regexp/tree/a468a5b)
- github.com/gregjones/httpcache: [901d907](https://github.com/gregjones/httpcache/tree/901d907)
- github.com/josharian/intern: [v1.0.0](https://github.com/josharian/intern/tree/v1.0.0)
- github.com/kisielk/errcheck: [v1.5.0](https://github.com/kisielk/errcheck/tree/v1.5.0)
- github.com/kisielk/gotool: [v1.0.0](https://github.com/kisielk/gotool/tree/v1.0.0)
- github.com/mailru/easyjson: [v0.9.0](https://github.com/mailru/easyjson/tree/v0.9.0)
- github.com/onsi/ginkgo/v2: [v2.21.0](https://github.com/onsi/ginkgo/tree/v2.21.0)
- github.com/onsi/gomega: [v1.35.1](https://github.com/onsi/gomega/tree/v1.35.1)
- github.com/pkg/errors: [v0.9.1](https://github.com/pkg/errors/tree/v0.9.1)
- github.com/yuin/goldmark: [v1.4.13](https://github.com/yuin/goldmark/tree/v1.4.13)
- github.com/zeebo/errs: [v1.4.0](https://github.com/zeebo/errs/tree/v1.4.0)
- go.uber.org/atomic: v1.11.0
- golang.org/x/telemetry: 1a19826
- golang.org/x/xerrors: 5ec99f8
- gopkg.in/yaml.v2: v2.4.0
- sigs.k8s.io/structured-merge-diff/v4: v4.6.0
