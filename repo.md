# Go Repo Layout — IDP Shared Library & Multi-CLI Monorepo

## Recommended Layout

```
idp/
├── cmd/                          # Binary entrypoints — one dir per CLI binary
│   ├── idp/                      # Main IDP CLI (scaffolding, registry ops)
│   │   └── main.go
│   ├── scaffold/                 # Standalone scaffolding binary
│   │   └── main.go
│   └── registry/                 # Registry query/management CLI
│       └── main.go
│
├── pkg/                          # Public, importable integration packages
│   ├── github/                   # GitHub API integration
│   │   ├── client.go
│   │   ├── repos.go
│   │   ├── teams.go
│   │   └── github_test.go
│   ├── circleci/                 # CircleCI integration
│   │   ├── client.go
│   │   ├── pipelines.go
│   │   └── circleci_test.go
│   ├── argocd/                   # ArgoCD integration
│   │   ├── client.go
│   │   ├── applications.go
│   │   └── argocd_test.go
│   ├── jira/                     # Jira integration
│   │   ├── client.go
│   │   ├── components.go
│   │   └── jira_test.go
│   ├── sonarcloud/               # SonarCloud integration
│   │   ├── client.go
│   │   ├── projects.go
│   │   └── sonarcloud_test.go
│   ├── port/                     # Port.io Developer Portal integration
│   │   ├── client.go
│   │   ├── entities.go
│   │   ├── blueprints.go
│   │   └── port_test.go
│   └── registry/                 # Central Components Registry client
│       ├── registry.go
│       ├── types.go
│       └── registry_test.go
│
├── internal/                     # Private packages — not importable externally
│   ├── config/                   # Config loading (env, YAML, flags)
│   │   ├── config.go
│   │   └── config_test.go
│   ├── scaffold/                 # Scaffolding orchestration logic
│   │   ├── orchestrator.go
│   │   ├── templates.go
│   │   └── scaffold_test.go
│   ├── validator/                # Shared input/entity validation
│   │   ├── validator.go
│   │   └── validator_test.go
│   └── httpclient/               # Shared base HTTP client (retries, auth, logging)
│       ├── client.go
│       └── client_test.go
│
├── api/                          # Shared types, interfaces, contracts
│   ├── types.go                  # Common domain types (Component, Service, Owner)
│   └── interfaces.go             # Integration interfaces (Scaffolder, Registrar, etc.)
│
├── templates/                    # Scaffolding templates (embedded via go:embed)
│   ├── github/
│   │   ├── repo-config.yaml.tmpl
│   │   └── codeowners.tmpl
│   ├── circleci/
│   │   └── config.yml.tmpl
│   └── service/
│       ├── java-springboot/
│       ├── dotnet/
│       ├── kotlin/
│       ├── nodejs/
│       └── golang/
│
├── scripts/                      # Build, release, and dev helper scripts
│   ├── build.sh
│   └── generate.sh
│
├── docs/                         # Additional documentation
│
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Key Design Decisions

### `cmd/` — Multiple CLI Binaries
Each subdirectory in `cmd/` produces a separate binary. They share packages from `pkg/` and `internal/` but have independent entrypoints. Build all with:

```bash
go build ./cmd/...
```

Or individually:
```bash
go build -o bin/idp ./cmd/idp
go build -o bin/scaffold ./cmd/scaffold
```

### `pkg/` — Public Integration Packages
Importable by external consumers (other repos, agents, automation). Each integration package should:
- Define its own `Client` struct with constructor accepting options/config
- Expose clean, domain-oriented methods (not raw HTTP wrappers)
- Be independently testable with interfaces for mocking

```go
// pkg/github/client.go
type Client struct { ... }

func New(token string, opts ...Option) *Client { ... }
func (c *Client) CreateRepo(ctx context.Context, opts CreateRepoOptions) (*Repo, error) { ... }
```

### `internal/` — Orchestration & Shared Plumbing
Business logic that ties integrations together lives here. The `scaffold` orchestrator calls `pkg/github`, `pkg/circleci`, `pkg/jira`, etc. in sequence. Not exported — keeps the public API surface clean.

### `api/` — Shared Domain Types & Interfaces
Defining interfaces here avoids circular imports between `pkg/` packages and `internal/`:

```go
// api/interfaces.go
type Scaffolder interface {
    Scaffold(ctx context.Context, spec ComponentSpec) error
}

type Registrar interface {
    Register(ctx context.Context, component Component) error
    Exists(ctx context.Context, id string) (bool, error)
}
```

This also makes it trivial for AI agents to discover and call capabilities — they can inspect `api/interfaces.go` as the authoritative contract.

### `templates/` — Embedded Templates
Use `go:embed` to bundle scaffolding templates directly into binaries:

```go
//go:embed templates
var TemplateFS embed.FS
```

No runtime file path dependencies — templates ship with the binary.

---

## Module Strategy

Single `go.mod` at the root — all packages share one module. This is the right call for a cohesive IDP toolset where integrations are tightly coupled.

```
module github.com/kubra-hq/idp

go 1.22
```

Only split into multiple modules if packages genuinely need independent versioning and release cadences.

---

## Makefile Targets

```makefile
.PHONY: build test lint tidy

build:
	go build -o bin/ ./cmd/...

test:
	go test ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

generate:
	go generate ./...
```

---

## Agentic Workflow Friendliness

The `api/interfaces.go` file acts as a machine-readable manifest of platform capabilities. An AI agent can:
- Call `registry.Exists(ctx, id)` before scaffolding to prevent duplicates
- Invoke `scaffold.Orchestrate(ctx, spec)` from natural language instructions
- Query `registry.List(ctx, filter)` for service metadata lookups

This layout makes the IDP repo a natural **tool provider** in agentic pipelines.