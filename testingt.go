package testcmp

// TestingT is the interface accepted by Equal and NotEqual.
// It is compatible with *testing.T, *testing.B, and testify suite types.
type TestingT interface {
	Errorf(format string, args ...interface{})
	Helper()
}
