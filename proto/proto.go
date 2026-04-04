// Package proto provides a cmp.Option helper for comparing protobuf messages
// with go-cmp.
//
// Example:
//
//	testcmp.Equal(t, want, got, proto.Transform())
package proto

import (
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

// Transform returns an option that makes go-cmp safe to use with protobuf structs.
// It uses protocmp.Transform(), which compares proto messages by their logical value
// rather than their internal unexported fields.
//
// Always use this when comparing proto-generated types. Without it, go-cmp panics
// on unexported fields injected by the proto runtime.
//
// Example:
//
//	testcmp.Equal(t, want, got, proto.Transform())
func Transform() cmp.Option {
	return protocmp.Transform()
}
