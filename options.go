package testcmp

import (
	"time"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// IgnoreTimestamps returns an option that ignores all time.Time values during comparison.
// Useful when structs contain CreatedAt, UpdatedAt, or similar fields set by the database.
//
// This ignores time.Time by type, so it applies to all time.Time fields and values,
// not just those with specific names. Use IgnoreAuditFields if you only want to ignore
// CreatedAt and UpdatedAt.
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
// Unlike IgnoreTimestamps, this targets fields by name, so other time.Time fields
// (e.g. ScheduledAt, BirthDate) are still compared.
//
// Example:
//
//	testcmp.Equal(t, want, got, testcmp.IgnoreAuditFields())
func IgnoreAuditFields() cmp.Option {
	return cmp.FilterPath(func(p cmp.Path) bool {
		sf, ok := p.Last().(cmp.StructField)
		if !ok {
			return false
		}
		switch sf.String() {
		case ".CreatedAt", ".UpdatedAt", ".CreatedBy", ".UpdatedBy":
			return true
		default:
			return false
		}
	}, cmp.Ignore())
}

// IgnoreIDs returns an option that ignores any struct field named exactly "ID"
// or ending with a lowercase letter followed by "ID", matching the standard Go
// naming convention (e.g. SupplierID, UserID, AssessmentID).
//
// Fields like VALID or SQUID are not matched because the character before "ID"
// is uppercase.
//
// If you need to ignore fields that don't follow this convention, use
// cmpopts.IgnoreFields directly.
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
		// sf.String() returns ".FieldName"
		name := sf.String()[1:] // strip leading dot
		if name == "ID" {
			return true
		}
		if len(name) < 3 {
			return false
		}
		// Match fields ending in "ID" where the preceding character is lowercase.
		return name[len(name)-2:] == "ID" && unicode.IsLower(rune(name[len(name)-3]))
	}, cmp.Ignore())
}
