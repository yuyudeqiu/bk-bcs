# AGENTS.md -- bcs-ingress-controller

## Project Snapshot

BCS Ingress Controller -- Kubernetes Operator managing network extension CRDs
(PortPool, PortBinding, Listener, Ingress, HostNetPortPool).
Go 1.20+, controller-runtime v0.6.3, go-restful HTTP API, Prometheus metrics.
go.mod lives at bcs-network/ parent; CRD types in ../../kubernetes/apis/networkextension/v1/.
Completed feature: HostNetPortPool hostNetwork port allocation (branch 001-hostnet-port-allocation).

## Build and Test

Build binary (run from bcs-ingress-controller/ dir):

    cd .. && make ingress-controller

Build + Docker image push (run from bcs-network/ dir):

    make ingress-controller
    VERSION=$(git describe --always)-$(date +%y.%m.%d)
    BUILD_DIR=../../../build/bcs.${VERSION}/bcs-runtime/bcs-k8s/bcs-network/bcs-ingress-controller
    docker build -t mirrors.tencent.com/<your_repository>/bcs-ingress-controller:latest -f ${BUILD_DIR}/Dockerfile ${BUILD_DIR}
    docker push mirrors.tencent.com/<your_repository>/bcs-ingress-controller:latest

Unit tests with coverage report:

    cd .. && make test-ingress-controller

Run single package tests:

    cd .. && go test -v -run TestReconcile ./bcs-ingress-controller/hostnetportcontroller/...

## K8s Cluster Deploy

    export KUBECONFIG=<path_to_your_kubeconfig>
    kubectl rollout restart -n bcs-system deployment/bcsingresscontroller

## Code Conventions

- Code comments: **English**. Docs / PR / chat: **Chinese**
- Format: gofmt / goimports, explicit error handling, no naked returns
- Constants in internal/constant/constant.go -- never use string literals for annotation keys
- Controllers: controller-runtime Reconcile, must be **idempotent**
- Tests: **table-driven**, colocated *_test.go, fake client from controller-runtime/pkg/client/fake
- Metrics: init() registration in internal/metrics/*.go, namespace bkbcs_ingressctrl
- Logging: bcs-common/common/blog -- not stdlib log or klog
- HTTP: go-restful WebService, routes registered in internal/httpsvr/httpserver.go InitRouters()

## Patterns -- DO

**New controller** -- follow portpoolcontroller/portpool_controller.go:
- Struct: ctx, client.Client, domain cache, record.EventRecorder
- Constructor NewXxxReconciler(ctx, cli, cache, eventer)
- SetupWithManager(mgr) with For(primaryCRD) + Watches(source.Kind)
- Reconcile(): fetch resource -> handle IsNotFound gracefully -> business logic -> update status
- Register in main.go via SetupWithManager(mgr)

**New metrics** -- follow internal/metrics/portpool.go:
- Package-level var block with prometheus.NewGaugeVec / CounterVec / HistogramVec
- init() calls metrics.Registry.MustRegister(...)
- Exported helpers: ReportXxx(...), CleanXxx(...), IncreaseXxx(...)

**New HTTP endpoint** -- follow internal/httpsvr/portpool.go:
- Handler method (h *HttpServerClient) handlerName(req, resp)
- Register in InitRouters() at internal/httpsvr/httpserver.go

**New cache** -- follow internal/portpoolcache/:
- Thread-safe with sync.RWMutex, types in separate types.go
- RebuildFromAPIServer() for cold-start / leader election recovery

**New constants** -- append to internal/constant/constant.go with exported const + godoc comment

## Patterns -- DON'T

- Import log / klog directly -- use blog
- Hardcode annotation keys -- add constant in internal/constant/constant.go
- Skip error checking on k8s API calls (client.Get, client.Update, etc.)
- Create controllers without registering in main.go
- Assume go.mod is in this directory -- it is at ../go.mod, run go mod tidy from ../

## Key Files

| Purpose | Path |
|---------|------|
| Entry point, wires all controllers / checkers / HTTP | main.go |
| All constants and annotation keys | internal/constant/constant.go |
| CLI flags and controller config | internal/option/ |
| HTTP API route registration | internal/httpsvr/httpserver.go |
| Webhook admission handlers | internal/webhookserver/ |
| CRD Go type definitions | ../../kubernetes/apis/networkextension/v1/ |
| CRD YAML manifests | ../../kubernetes/config/crd/bases/ |
| Generated clientset / listers / informers | ../../kubernetes/generated/ |
| Build Makefile | ../Makefile |

## Controllers and CRDs

| Controller | CRD / Resource | Directory |
|------------|---------------|-----------|
| Ingress | networkextension.Ingress | ingresscontroller/ |
| Listener | networkextension.Listener | listenercontroller/ |
| PortPool | networkextension.PortPool | portpoolcontroller/ |
| PortBinding | networkextension.PortBinding | portbindingcontroller/ |
| HostNetPortPool | networkextension.HostNetPortPool | hostnetportcontroller/ |
| Namespace | core/v1.Namespace | namespacecontroller/ |
| Node | core/v1.Node | nodecontroller/ |

## Module Layout

    bcs-ingress-controller/
      main.go                     # Wires everything
      {name}controller/           # One dir per controller
        controller.go             # Reconciler + SetupWithManager
        *_test.go                 # Colocated tests
      internal/
        constant/                 # Shared constants
        metrics/                  # Prometheus metrics (one file per subsystem)
        httpsvr/                  # REST API handlers + route registration
        check/                    # Periodic consistency checkers
        cloud/                    # Cloud adapters (aws/ azure/ gcp/ tencentcloud/)
        webhookserver/            # Admission webhooks
        hostnetportpoolcache/     # HostNetPortPool in-memory cache
        portpoolcache/            # PortPool in-memory cache
        generator/                # Ingress to Listener conversion
        ingresscache/             # Service/workload cache for Ingress
        nodecache/                # Node metadata cache
        apiclient/                # External API client helpers
        option/                   # CLI option parsing
        worker/                   # Worker/sync utilities
      bcs-ingress-inspector/      # Separate diagnostic binary (NOT main controller)
      specs/                      # Feature design documents
      benchmark/                  # Performance test scripts and fixtures

## JIT Search Commands

Find all Reconcile implementations:

    rg -n "func.*Reconcile\(" --type go

Find annotation / constant definitions:

    rg -n "const\b" internal/constant/constant.go

Find Prometheus metric definitions:

    rg -n "prometheus\.New" internal/metrics/

Find HTTP routes:

    rg -n "ws\.Route" internal/httpsvr/

Find CRD type structs:

    rg -n "type .*(Spec|Status) struct" ../../kubernetes/apis/networkextension/v1/

Find where controllers are registered:

    rg -n "SetupWithManager" main.go

Find webhook handlers:

    rg -n "Handle\(" internal/webhookserver/

Find all test files:

    rg -l "_test\.go" --type go

## Pre-PR Checklist

Run tests before submitting:

    cd .. && make test-ingress-controller && echo "All tests OK"

Verify:
- New constants added to internal/constant/constant.go
- New controllers registered in main.go via SetupWithManager
- New metrics registered via init() in internal/metrics/
- New HTTP routes added in internal/httpsvr/httpserver.go InitRouters()
- CRD type changes: regenerate deepcopy and manifests in ../../kubernetes/
- gofmt / goimports clean

## Common Gotchas

- go.mod is at **bcs-network/** level -- always run go mod tidy from ../
- CRD types live in a **separate Go module** (../../kubernetes/), linked via replace directive
- controller-runtime v0.6.3 + effective k8s client v0.18.6 -- old API style, no generics
- bcs-ingress-inspector/ is a **separate binary** with its own main.go, not part of this controller
- The cloud/ adapters have vendor-specific SDKs -- only import the one you need

## Completed Features

- ~~001-hostnet-port-allocation~~: HostNetPortPool dynamic port allocation for hostNetwork Pods (done)
- Design docs: specs/001-hostnet-port-allocation/
