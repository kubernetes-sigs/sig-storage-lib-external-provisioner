# Release notes for v12.0.0

# Changelog since v11.1.0

## Changes by Kind

### Urgent Upgrade Notes

- Signatures of two functions changed to ask for string-typed rate limiters.
  - RateLimiter
  - CreateProvisionedPVLimiter ([#191](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/191), [@grant-he](https://github.com/grant-he))

### Bug or Regression

- Avoid workqueue depth metric goes negative ([#190](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/190), [@huww98](https://github.com/huww98))

### Uncategorized

- Update to Kubernetes 1.34 ([#193](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/193), [@rhrmo](https://github.com/rhrmo))

## Dependencies

### Added
- github.com/envoyproxy/go-control-plane/envoy: [v1.32.4](https://github.com/envoyproxy/go-control-plane/tree/envoy/v1.32.4)
- github.com/envoyproxy/go-control-plane/ratelimit: [v0.1.0](https://github.com/envoyproxy/go-control-plane/tree/ratelimit/v0.1.0)
- github.com/go-jose/go-jose/v4: [v4.1.1](https://github.com/go-jose/go-jose/tree/v4.1.1)
- github.com/go-openapi/swag/cmdutils: [v0.24.0](https://github.com/go-openapi/swag/tree/cmdutils/v0.24.0)
- github.com/go-openapi/swag/conv: [v0.24.0](https://github.com/go-openapi/swag/tree/conv/v0.24.0)
- github.com/go-openapi/swag/fileutils: [v0.24.0](https://github.com/go-openapi/swag/tree/fileutils/v0.24.0)
- github.com/go-openapi/swag/jsonname: [v0.24.0](https://github.com/go-openapi/swag/tree/jsonname/v0.24.0)
- github.com/go-openapi/swag/jsonutils: [v0.24.0](https://github.com/go-openapi/swag/tree/jsonutils/v0.24.0)
- github.com/go-openapi/swag/loading: [v0.24.0](https://github.com/go-openapi/swag/tree/loading/v0.24.0)
- github.com/go-openapi/swag/mangling: [v0.24.0](https://github.com/go-openapi/swag/tree/mangling/v0.24.0)
- github.com/go-openapi/swag/netutils: [v0.24.0](https://github.com/go-openapi/swag/tree/netutils/v0.24.0)
- github.com/go-openapi/swag/stringutils: [v0.24.0](https://github.com/go-openapi/swag/tree/stringutils/v0.24.0)
- github.com/go-openapi/swag/typeutils: [v0.24.0](https://github.com/go-openapi/swag/tree/typeutils/v0.24.0)
- github.com/go-openapi/swag/yamlutils: [v0.24.0](https://github.com/go-openapi/swag/tree/yamlutils/v0.24.0)
- github.com/grafana/regexp: [a468a5b](https://github.com/grafana/regexp/tree/a468a5b)
- github.com/spiffe/go-spiffe/v2: [v2.5.0](https://github.com/spiffe/go-spiffe/tree/v2.5.0)
- github.com/zeebo/errs: [v1.4.0](https://github.com/zeebo/errs/tree/v1.4.0)
- go.uber.org/atomic: v1.11.0
- go.yaml.in/yaml/v2: v2.4.2
- go.yaml.in/yaml/v3: v3.0.4
- golang.org/x/tools/go/expect: v0.1.0-deprecated
- golang.org/x/tools/go/packages/packagestest: v0.1.1-deprecated
- gonum.org/v1/gonum: v0.16.0
- sigs.k8s.io/structured-merge-diff/v6: v6.3.0

### Changed
- cel.dev/expr: v0.16.2 → v0.24.0
- cloud.google.com/go/compute/metadata: v0.5.2 → v0.7.0
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.24.2 → v1.29.0](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/compare/detectors/gcp/v1.24.2...detectors/gcp/v1.29.0)
- github.com/cncf/xds/go: [b4127c9 → 2ac532f](https://github.com/cncf/xds/compare/b4127c9...2ac532f)
- github.com/emicklei/go-restful/v3: [v3.12.2 → v3.13.0](https://github.com/emicklei/go-restful/compare/v3.12.2...v3.13.0)
- github.com/envoyproxy/go-control-plane: [v0.13.1 → v0.13.4](https://github.com/envoyproxy/go-control-plane/compare/v0.13.1...v0.13.4)
- github.com/envoyproxy/protoc-gen-validate: [v1.1.0 → v1.2.1](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.1.0...v1.2.1)
- github.com/fxamacker/cbor/v2: [v2.8.0 → v2.9.0](https://github.com/fxamacker/cbor/compare/v2.8.0...v2.9.0)
- github.com/go-openapi/jsonpointer: [v0.21.1 → v0.22.0](https://github.com/go-openapi/jsonpointer/compare/v0.21.1...v0.22.0)
- github.com/go-openapi/jsonreference: [v0.21.0 → v0.21.1](https://github.com/go-openapi/jsonreference/compare/v0.21.0...v0.21.1)
- github.com/go-openapi/swag: [v0.23.1 → v0.24.1](https://github.com/go-openapi/swag/compare/v0.23.1...v0.24.1)
- github.com/golang/glog: [v1.2.2 → v1.2.5](https://github.com/golang/glog/compare/v1.2.2...v1.2.5)
- github.com/google/gnostic-models: [v0.6.9 → v0.7.0](https://github.com/google/gnostic-models/compare/v0.6.9...v0.7.0)
- github.com/miekg/dns: [v1.1.66 → v1.1.68](https://github.com/miekg/dns/compare/v1.1.66...v1.1.68)
- github.com/modern-go/reflect2: [v1.0.2 → 35a7c28](https://github.com/modern-go/reflect2/compare/v1.0.2...35a7c28)
- github.com/prometheus/client_golang: [v1.22.0 → v1.23.0](https://github.com/prometheus/client_golang/compare/v1.22.0...v1.23.0)
- github.com/prometheus/common: [v0.64.0 → v0.66.0](https://github.com/prometheus/common/compare/v0.64.0...v0.66.0)
- github.com/prometheus/procfs: [v0.16.1 → v0.17.0](https://github.com/prometheus/procfs/compare/v0.16.1...v0.17.0)
- github.com/spf13/pflag: [v1.0.5 → v1.0.6](https://github.com/spf13/pflag/compare/v1.0.5...v1.0.6)
- github.com/stretchr/testify: [v1.10.0 → v1.11.1](https://github.com/stretchr/testify/compare/v1.10.0...v1.11.1)
- go.opentelemetry.io/contrib/detectors/gcp: v1.31.0 → v1.36.0
- go.opentelemetry.io/otel/metric: v1.33.0 → v1.37.0
- go.opentelemetry.io/otel/sdk/metric: v1.31.0 → v1.37.0
- go.opentelemetry.io/otel/sdk: v1.31.0 → v1.37.0
- go.opentelemetry.io/otel/trace: v1.33.0 → v1.37.0
- go.opentelemetry.io/otel: v1.33.0 → v1.37.0
- golang.org/x/crypto: v0.38.0 → v0.41.0
- golang.org/x/mod: v0.24.0 → v0.27.0
- golang.org/x/net: v0.40.0 → v0.43.0
- golang.org/x/sync: v0.14.0 → v0.16.0
- golang.org/x/sys: v0.33.0 → v0.35.0
- golang.org/x/telemetry: bda5523 → 1a19826
- golang.org/x/term: v0.32.0 → v0.34.0
- golang.org/x/text: v0.25.0 → v0.28.0
- golang.org/x/time: v0.11.0 → v0.12.0
- golang.org/x/tools: v0.33.0 → v0.36.0
- google.golang.org/genproto/googleapis/api: 796eee8 → 8d1bb00
- google.golang.org/genproto/googleapis/rpc: 9240e9c → ef028d9
- google.golang.org/grpc: v1.69.0 → v1.75.0
- google.golang.org/protobuf: v1.36.6 → v1.36.8
- gopkg.in/evanphx/json-patch.v4: v4.12.0 → v4.13.0
- k8s.io/api: v0.33.1 → v0.34.0
- k8s.io/apimachinery: v0.33.1 → v0.34.0
- k8s.io/client-go: v0.33.1 → v0.34.0
- k8s.io/gengo/v2: a7b603a → 85fd79d
- k8s.io/kube-openapi: c8a335a → 7fc2783
- k8s.io/utils: 0f33e8f → 0af2bda
- sigs.k8s.io/json: cfa47c3 → 2d32026
- sigs.k8s.io/structured-merge-diff/v4: v4.7.0 → v4.6.0
- sigs.k8s.io/yaml: v1.4.0 → v1.6.0

### Removed
- github.com/census-instrumentation/opencensus-proto: [v0.4.1](https://github.com/census-instrumentation/opencensus-proto/tree/v0.4.1)
