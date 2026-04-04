package testcmp

import (
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

// IgnoreTimestamps returns an option that ignores all time.Time fields during comparison.
// Useful when structs contain CreatedAt, UpdatedAt, or similar fields set by the database.
//
// Example:
//
//	testcmp.Equal(t, want, got, testcmp.IgnoreTimestamps())
func IgnoreTimestamps() cmp.Option {
	return cmpopts.IgnoreTypes(time.Time{})
}

// IgnoreAuditFields returns an option that ignores common audit fields by name:
// CreatedAt, UpdatedAt, CreatedBy, UpdatedBy.
//
// It combines IgnoreTimestamps with a field name filter for the *By fields.
// Use this when your domain structs embed audit metadata you do not control in tests.
//
// Example:
//
//	testcmp.Equal(t, want, got, testcmp.IgnoreAuditFields())
func IgnoreAuditFields() cmp.Option {
	return cmp.Options{
		IgnoreTimestamps(),
		cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
			switch k {
			case "CreatedBy", "UpdatedBy":
				return true
			default:
				return false
			}
		}),
		cmp.FilterPath(func(p cmp.Path) bool {
			sf, ok := p.Last().(cmp.StructField)
			if !ok {
				return false
			}
			switch sf.String() {
			case ".CreatedBy", ".UpdatedBy":
				return true
			default:
				return false
			}
		}, cmp.Ignore()),
	}
}

// IgnoreIDs returns an option that ignores any field whose name ends in "ID".
// This covers both uuid.UUID and int64 ID fields by field name, not by type,
// so it is safe to use when your structs have mixed ID types.
//
// Fields matched: ID, SupplierID, AssessmentID, UserID, etc.
//
// If you have non-ID fields that happen to end in "ID", use IgnoreFields instead.
//
// Example:
//
//	testcmp.Equal(t, want, got, testcmp.IgnoreIDs())
func IgnoreIDs() cmp.Option {
	return cmp.FilterPath(func(p cmp.Path) bool {
		sf, ok := p.Last().(cmp.StructField)
		if !ok {
			return false
		}
		name := sf.String()
		// sf.String() returns ".FieldName"
		return strings.HasSuffix(name, "ID")
	}, cmp.Ignore())
}

// IgnoreFields returns an option that ignores the named fields on the given struct type.
// This is a typed wrapper around cmpopts.IgnoreFields.
//
// Use this when IgnoreIDs is too broad and you need precise control.
//
// Example:
//
//	testcmp.Equal(t, want, got,
//	    testcmp.IgnoreFields(Supplier{}, "ID", "ExternalID"),
//	)
func IgnoreFields[T any](typ T, fields ...string) cmp.Option {
	return cmpopts.IgnoreFields(typ, fields...)
}

// Proto returns an option that makes go-cmp safe to use with protobuf structs.
// It uses protocmp.Transform(), which compares proto messages by their logical value
// rather than their internal unexported fields.
//
// Always use this when comparing proto-generated types. Without it, go-cmp panics
// on unexported fields injected by the proto runtime.
//
// Requires: google.golang.org/protobuf/testing/protocmp
//
// Example:
//
//	testcmp.Equal(t, want, got, testcmp.Proto())
func Proto() cmp.Option {
	return protocmp.Transform()
}
