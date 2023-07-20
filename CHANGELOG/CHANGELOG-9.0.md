# Release notes for v9.0.3

# Changelog since v9.0.2

## Changes by Kind

### Uncategorized

- Added debug logs to shouldDelete function ([#146](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/146), [@amacaskill](https://github.com/amacaskill))

# Release notes for v9.0.2

# Changelog since v9.0.1

## Changes by Kind

### Uncategorized

- Fix indefinite stuck Pending pod on a deleted node ([#139](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/139), [@sunnylovestiramisu](https://github.com/sunnylovestiramisu))

# Release notes for v9.0.1

# Changelog since v9.0.0

## Changes by Kind

### Uncategorized

- Fixed go.mod for v9.

# Release notes for v9.0.0

# Changelog since v8.0.0

## Changes by Kind

### Uncategorized

- Action Needed: Updated leader election to use Endpoints + Lease. All provisioners based on this library must have RBAC permissions to create/update Lease objects in coordination.k8s.io/v1 API. ([#120](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/120), [@NikhilSharmaWe](https://github.com/NikhilSharmaWe))
- Added a new label `source` to `controller_persistentvolumeclaim_provision_total` metric. ([#128](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/128), [@RaunakShah](https://github.com/RaunakShah))
- Added support for `external-provisioner.volume.kubernetes.io/finalizer` on statically provisioned volumes. ([#129](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/pull/129), [@deepakkinni](https://github.com/deepakkinni))


## Dependencies

### Added
- github.com/armon/go-socks5: [e753329](https://github.com/armon/go-socks5/tree/e753329)
- github.com/asaskevich/govalidator: [f61b66f](https://github.com/asaskevich/govalidator/tree/f61b66f)
- github.com/cncf/udpa/go: [269d4d4](https://github.com/cncf/udpa/go/tree/269d4d4)
- github.com/creack/pty: [v1.1.9](https://github.com/creack/pty/tree/v1.1.9)
- github.com/emicklei/go-restful/v3: [v3.9.0](https://github.com/emicklei/go-restful/v3/tree/v3.9.0)
- github.com/go-gl/glfw: [e6da0ac](https://github.com/go-gl/glfw/tree/e6da0ac)
- github.com/google/gnostic: [v0.5.7-v3refs](https://github.com/google/gnostic/tree/v0.5.7-v3refs)
- github.com/google/martian/v3: [v3.0.0](https://github.com/google/martian/v3/tree/v3.0.0)
- github.com/josharian/intern: [v1.0.0](https://github.com/josharian/intern/tree/v1.0.0)
- github.com/mitchellh/mapstructure: [v1.1.2](https://github.com/mitchellh/mapstructure/tree/v1.1.2)
- github.com/moby/spdystream: [v0.2.0](https://github.com/moby/spdystream/tree/v0.2.0)
- github.com/niemeyer/pretty: [a10e7ca](https://github.com/niemeyer/pretty/tree/a10e7ca)
- github.com/onsi/ginkgo/v2: [v2.4.0](https://github.com/onsi/ginkgo/v2/tree/v2.4.0)
- github.com/stoewer/go-strcase: [v1.2.0](https://github.com/stoewer/go-strcase/tree/v1.2.0)
- github.com/yuin/goldmark: [v1.2.1](https://github.com/yuin/goldmark/tree/v1.2.1)
- golang.org/x/term: v0.3.0
- gopkg.in/yaml.v3: v3.0.1
- rsc.io/quote/v3: v3.1.0
- rsc.io/sampler: v1.3.0
- sigs.k8s.io/json: f223a00

### Changed
- cloud.google.com/go/bigquery: v1.0.1 → v1.8.0
- cloud.google.com/go/datastore: v1.0.0 → v1.1.0
- cloud.google.com/go/pubsub: v1.0.1 → v1.3.1
- cloud.google.com/go/storage: v1.0.0 → v1.10.0
- cloud.google.com/go: v0.51.0 → v0.65.0
- github.com/PuerkitoBio/purell: [v1.0.0 → v1.1.1](https://github.com/PuerkitoBio/purell/compare/v1.0.0...v1.1.1)
- github.com/PuerkitoBio/urlesc: [5bd2802 → de5bf2a](https://github.com/PuerkitoBio/urlesc/compare/5bd2802...de5bf2a)
- github.com/envoyproxy/go-control-plane: [5f8ba28 → v0.9.4](https://github.com/envoyproxy/go-control-plane/compare/5f8ba28...v0.9.4)
- github.com/evanphx/json-patch: [v4.9.0+incompatible → v4.12.0+incompatible](https://github.com/evanphx/json-patch/compare/v4.9.0...v4.12.0)
- github.com/go-gl/glfw/v3.3/glfw: [12ad95a → 6f7a984](https://github.com/go-gl/glfw/v3.3/glfw/compare/12ad95a...6f7a984)
- github.com/go-logr/logr: [v0.2.0 → v1.2.3](https://github.com/go-logr/logr/compare/v0.2.0...v1.2.3)
- github.com/go-openapi/jsonpointer: [46af16f → v0.19.5](https://github.com/go-openapi/jsonpointer/compare/46af16f...v0.19.5)
- github.com/go-openapi/jsonreference: [13c6e35 → v0.20.0](https://github.com/go-openapi/jsonreference/compare/13c6e35...v0.20.0)
- github.com/go-openapi/swag: [1d0bd11 → v0.19.14](https://github.com/go-openapi/swag/compare/1d0bd11...v0.19.14)
- github.com/gogo/protobuf: [v1.3.1 → v1.3.2](https://github.com/gogo/protobuf/compare/v1.3.1...v1.3.2)
- github.com/golang/groupcache: [215e871 → 41bb18b](https://github.com/golang/groupcache/compare/215e871...41bb18b)
- github.com/golang/mock: [v1.3.1 → v1.4.4](https://github.com/golang/mock/compare/v1.3.1...v1.4.4)
- github.com/golang/protobuf: [v1.4.2 → v1.5.2](https://github.com/golang/protobuf/compare/v1.4.2...v1.5.2)
- github.com/google/btree: [v1.0.0 → v1.0.1](https://github.com/google/btree/compare/v1.0.0...v1.0.1)
- github.com/google/go-cmp: [v0.4.0 → v0.5.9](https://github.com/google/go-cmp/compare/v0.4.0...v0.5.9)
- github.com/google/pprof: [d4f498a → 1a94d86](https://github.com/google/pprof/compare/d4f498a...1a94d86)
- github.com/google/uuid: [v1.1.1 → v1.1.2](https://github.com/google/uuid/compare/v1.1.1...v1.1.2)
- github.com/imdario/mergo: [v0.3.5 → v0.3.6](https://github.com/imdario/mergo/compare/v0.3.5...v0.3.6)
- github.com/json-iterator/go: [v1.1.10 → v1.1.12](https://github.com/json-iterator/go/compare/v1.1.10...v1.1.12)
- github.com/kisielk/errcheck: [v1.2.0 → v1.5.0](https://github.com/kisielk/errcheck/compare/v1.2.0...v1.5.0)
- github.com/kr/text: [v0.1.0 → v0.2.0](https://github.com/kr/text/compare/v0.1.0...v0.2.0)
- github.com/mailru/easyjson: [d5b7844 → v0.7.6](https://github.com/mailru/easyjson/compare/d5b7844...v0.7.6)
- github.com/modern-go/reflect2: [v1.0.1 → v1.0.2](https://github.com/modern-go/reflect2/compare/v1.0.1...v1.0.2)
- github.com/munnerz/goautoneg: [a547fc6 → a7dc8b6](https://github.com/munnerz/goautoneg/compare/a547fc6...a7dc8b6)
- github.com/onsi/gomega: [v1.9.0 → v1.23.0](https://github.com/onsi/gomega/compare/v1.9.0...v1.23.0)
- github.com/stretchr/testify: [v1.4.0 → v1.8.0](https://github.com/stretchr/testify/compare/v1.4.0...v1.8.0)
- go.opencensus.io: v0.22.2 → v0.22.4
- golang.org/x/crypto: 75b2880 → v0.1.0
- golang.org/x/exp: da58074 → 6cc2880
- golang.org/x/lint: fdd1cda → 738671d
- golang.org/x/mod: c90efee → 86c51ed
- golang.org/x/net: ab34263 → 1e63c2f
- golang.org/x/oauth2: bf48bf1 → ee48083
- golang.org/x/sync: cd5d95a → 886fb93
- golang.org/x/sys: ed371f2 → v0.3.0
- golang.org/x/text: v0.3.3 → v0.5.0
- golang.org/x/time: 555d28b → 90d013b
- golang.org/x/tools: 7b8e75d → v0.1.12
- golang.org/x/xerrors: 9bdfabe → 5ec99f8
- google.golang.org/api: v0.15.0 → v0.30.0
- google.golang.org/appengine: v1.6.5 → v1.6.7
- google.golang.org/genproto: cb27e3a → 1ed22bb
- google.golang.org/grpc: v1.27.0 → v1.31.0
- google.golang.org/protobuf: v1.24.0 → v1.28.1
- gopkg.in/check.v1: 41f04d3 → 8fa4692
- gopkg.in/yaml.v2: v2.2.8 → v2.4.0
- honnef.co/go/tools: v0.0.1-2019.2.3 → v0.0.1-2020.1.4
- k8s.io/api: v0.19.1 → v0.26.0
- k8s.io/apimachinery: v0.19.1 → v0.26.0
- k8s.io/client-go: v0.19.1 → v0.26.0
- k8s.io/gengo: 3a45101 → 485abfe
- k8s.io/klog/v2: v2.3.0 → v2.80.1
- k8s.io/kube-openapi: 6aeccd4 → 172d655
- k8s.io/utils: d5654de → 1a15be2
- sigs.k8s.io/structured-merge-diff/v4: v4.0.1 → v4.2.3
- sigs.k8s.io/yaml: v1.2.0 → v1.3.0

### Removed
- github.com/Azure/go-autorest/autorest/adal: [v0.8.2](https://github.com/Azure/go-autorest/autorest/adal/tree/v0.8.2)
- github.com/Azure/go-autorest/autorest/date: [v0.2.0](https://github.com/Azure/go-autorest/autorest/date/tree/v0.2.0)
- github.com/Azure/go-autorest/autorest/mocks: [v0.3.0](https://github.com/Azure/go-autorest/autorest/mocks/tree/v0.3.0)
- github.com/Azure/go-autorest/autorest: [v0.9.6](https://github.com/Azure/go-autorest/autorest/tree/v0.9.6)
- github.com/Azure/go-autorest/logger: [v0.1.0](https://github.com/Azure/go-autorest/logger/tree/v0.1.0)
- github.com/Azure/go-autorest/tracing: [v0.5.0](https://github.com/Azure/go-autorest/tracing/tree/v0.5.0)
- github.com/dgrijalva/jwt-go: [v3.2.0+incompatible](https://github.com/dgrijalva/jwt-go/tree/v3.2.0)
- github.com/docker/spdystream: [449fdfc](https://github.com/docker/spdystream/tree/449fdfc)
- github.com/emicklei/go-restful: [ff4f55a](https://github.com/emicklei/go-restful/tree/ff4f55a)
- github.com/fsnotify/fsnotify: [v1.4.9](https://github.com/fsnotify/fsnotify/tree/v1.4.9)
- github.com/ghodss/yaml: [73d445a](https://github.com/ghodss/yaml/tree/73d445a)
- github.com/go-openapi/spec: [6aced65](https://github.com/go-openapi/spec/tree/6aced65)
- github.com/googleapis/gnostic: [v0.4.1](https://github.com/googleapis/gnostic/tree/v0.4.1)
- github.com/hpcloud/tail: [v1.0.0](https://github.com/hpcloud/tail/tree/v1.0.0)
- github.com/onsi/ginkgo: [v1.12.0](https://github.com/onsi/ginkgo/tree/v1.12.0)
- github.com/spf13/afero: [v1.2.2](https://github.com/spf13/afero/tree/v1.2.2)
- gopkg.in/fsnotify.v1: v1.4.7
- gopkg.in/tomb.v1: dd63297
