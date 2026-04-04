# testcmp

A thin Go library that bridges [testify](https://github.com/stretchr/testify) and
[go-cmp](https://github.com/google/go-cmp). It provides assertion helpers that use
`cmp.Diff` under the hood, giving you structured, readable diffs on failure while
keeping the testify style you already use.

## Install

```bash
go get github.com/rnaveiras/testcmp
```

Requires **Go 1.24** or later.

## Quick start

```go
import "github.com/rnaveiras/testcmp"

func TestGetSupplier(t *testing.T) {
    want := Supplier{Name: "Acme Corp", Status: StatusActive}
    got  := store.GetSupplier(ctx, "acme")

    testcmp.Equal(t, want, got)
}
```

On failure you get a clear diff:

```
Error:    mismatch (-want +got):
          - Name: "Alice"
          + Name: "Bob"
```

## Option helpers

All helpers return a `cmp.Option` and live in the same package, so a single import covers
both assertions and options.

### IgnoreTimestamps

Ignores all `time.Time` fields during comparison.

```go
testcmp.Equal(t, want, got, testcmp.IgnoreTimestamps())
```

### IgnoreAuditFields

Ignores `CreatedAt`, `UpdatedAt`, `CreatedBy`, and `UpdatedBy`.

```go
testcmp.Equal(t, want, got, testcmp.IgnoreAuditFields())
```

### IgnoreIDs

Ignores struct fields named exactly `"ID"` or ending with a lowercase letter followed by
`"ID"`, matching the standard Go naming convention (e.g. `ID`, `SupplierID`, `UserID`).
Fields like `VALID` or `SQUID` are not matched.

```go
testcmp.Equal(t, want, got, testcmp.IgnoreIDs())
```

### Proto (sub-package)

The `testcmp/proto` sub-package provides protobuf support via `protocmp.Transform()`,
keeping `google.golang.org/protobuf` out of the core module's dependency tree.

```go
import "github.com/rnaveiras/testcmp/proto"

testcmp.Equal(t, want, got, proto.Transform())
```

### Combining options

Options compose naturally:

```go
testcmp.Equal(t, want, got,
    testcmp.IgnoreAuditFields(),
    testcmp.IgnoreIDs(),
    proto.Transform(),
)
```

## API

### Assertions

```go
func Equal[T any](t testcmp.TestingT, want, got T, opts ...cmp.Option) bool
func NotEqual[T any](t testcmp.TestingT, want, got T, opts ...cmp.Option) bool
```

`TestingT` requires `Errorf` and `Helper`, satisfied by `*testing.T`, `*testing.B`,
and testify suite types.

Both return `bool` matching testify's convention, so you can gate further assertions:

```go
if testcmp.Equal(t, want, got) {
    // safe to assert on got.Field
}
```

## Dependencies

| Dependency | Reason |
|---|---|
| `github.com/google/go-cmp` | Diff and comparison engine |
| `google.golang.org/protobuf` | `protocmp.Transform()` (only in `testcmp/proto` sub-module) |

## License

[MIT](LICENSE)
