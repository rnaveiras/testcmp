package testcmp_test

import (
	"fmt"
	"testing"
)

// tbRecorder captures Errorf calls for test assertions.
// it embeds a real testing.TB to satisfy the unexported method.
type tbRecorder struct {
	testing.TB
	failed bool
	errors []string
}

func newTBRecorder(t *testing.T) *tbRecorder {
	t.Helper()

	return &tbRecorder{TB: t}
}

func (r *tbRecorder) Errorf(format string, args ...any) {
	r.failed = true
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

func (r *tbRecorder) Helper() {}
