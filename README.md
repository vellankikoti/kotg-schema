# kotg-schema

Versioned wire contract between [kubilitics](https://github.com/vellankikoti/kubilitics)
core and [kotg.ai](https://github.com/vellankikoti/kotg.ai). One
protobuf schema, one Go module, two consumers, independent release
cadence on either side.

This repo's only job is to define the bytes that travel over the wire.
Everything else (servers, clients, business logic) lives in the
consumers. See the [integration design](https://github.com/vellankikoti/kubilitics/blob/main/docs/superpowers/specs/2026-04-18-ai-integration-design.md)
and [schema design](https://github.com/vellankikoti/kubilitics/blob/main/docs/superpowers/specs/2026-04-18-kotg-schema-design.md)
for the full picture.

## Wire contract

Three proto packages under `proto/kotg/v1/`:

| Package         | Services                                                        |
|-----------------|-----------------------------------------------------------------|
| `cluster.proto` | `ClusterRead`, `ClusterAction`, `ActionTemplate`                |
| `chat.proto`    | `Chat`, `AIControl`                                             |
| `common.proto`  | shared types (`ActionTier`, `ResourceRef`, `Diff`, `AuditExtras`, `HealthStatus`) |

Generated Go is committed under `gen/go/kotg/v1/` so consumers do not
need protoc.

## Five rules baked into the contract

1. **User identity travels in gRPC metadata** under key `kotg-user-token`. Never in a body field.
2. **Every mutating `ActionResult` carries `undo_token` + `undo_ttl_seconds` as REQUIRED fields.** Powers one-click rollback.
3. **`Citation` events SHOULD accompany any assistant claim about cluster state.** Powers the receipts UX. Treat missing-citation as a quality bug.
4. **`ActionPending` carries the full `Diff` and `tier`** so approval modals render without an extra round-trip.
5. **mTLS is mandatory.** No `tls_disabled` field. No plaintext fallback.

## Versioning

Semver, `v1.x.y`. Patch + minor are additive only — old generated code keeps compiling. Major version bumps require a coordinated release with a 3-month overlap window where both v1 and v2 are served. CI runs `buf breaking` against the previous tag on every PR.

## Using from Go

```go
import kotgv1 "github.com/kubilitics/kotg-schema/gen/go/kotg/v1"
```

See `examples/stub_cluster_server/` and `examples/stub_chat_server/` for minimal working servers.

## Development

```bash
brew install bufbuild/buf/buf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

buf lint
buf generate
go build ./...
go test ./...
```
