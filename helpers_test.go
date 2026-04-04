package testcmp_test

import "fmt"

// mockT implements assert.TestingT for testing assertion functions
// without using a real *testing.T.
type mockT struct {
	failed   bool
	messages []string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}

func (m *mockT) Helper() {}
