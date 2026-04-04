// Package testcmp bridges testify and go-cmp, providing assertion helpers
// that use cmp.Diff under the hood for structured, readable diffs on failure
// while keeping the testify style.
//
// It also ships pre-built cmp.Option helpers covering the most common cases:
// timestamps, IDs, and audit fields.
//
// Example:
//
//	testcmp.Equal(t, want, got, testcmp.IgnoreTimestamps())
package testcmp

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// Equal asserts that want and got are equal using cmp.Diff.
// opts are passed directly to cmp.Diff.
//
// Example:
//
//	testcmp.Equal(t, want, got, testcmp.IgnoreTimestamps())
func Equal[T any](t assert.TestingT, want, got T, opts ...cmp.Option) bool {
	if h, ok := t.(interface{ Helper() }); ok {
		h.Helper()
	}

	diff := cmp.Diff(want, got, opts...)
	if diff == "" {
		return true
	}

	t.Errorf("mismatch (-want +got):\n%s", diff)
	return false
}

// NotEqual asserts that want and got are not equal using cmp.Diff.
//
// Example:
//
//	testcmp.NotEqual(t, original, updated, testcmp.IgnoreTimestamps())
func NotEqual[T any](t assert.TestingT, want, got T, opts ...cmp.Option) bool {
	if h, ok := t.(interface{ Helper() }); ok {
		h.Helper()
	}

	diff := cmp.Diff(want, got, opts...)
	if diff != "" {
		return true
	}

	t.Errorf("values are equal, but should differ")
	return false
}
