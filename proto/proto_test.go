package proto_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rnaveiras/testcmp/proto"
)

func TestTransform_ReturnsCmpOption(t *testing.T) {
	opt := proto.Transform()
	require.NotNil(t, opt)

	// Verify it's usable as a cmp.Option by passing it to cmp.Diff
	// with simple non-proto values (should not panic).
	assert.NotPanics(t, func() {
		_ = cmp.Diff("a", "a", opt)
	})
}
