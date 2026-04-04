package proto_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rnaveiras/testcmp/proto"
)

func TestTransform_ReturnsCmpOption(t *testing.T) {
	opt := proto.Transform()
	if opt == nil {
		t.Fatal("Transform() should return a non-nil cmp.Option")
	}

	// Verify it's usable as a cmp.Option by passing it to cmp.Diff
	// with simple non-proto values (should not panic).
	_ = cmp.Diff("a", "a", opt)
}
