package testcmp_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rnaveiras/testcmp"
)

type mockT struct {
	failed   bool
	messages []string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.messages = append(m.messages, format)
}

func (m *mockT) Helper() {}

func TestEqual_SameValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.Equal(mt, S{Name: "Alice"}, S{Name: "Alice"})
	if !result {
		t.Error("Equal should return true for identical values")
	}
	if mt.failed {
		t.Error("Equal should not call Errorf for identical values")
	}
}

func TestEqual_DifferentValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.Equal(mt, S{Name: "Alice"}, S{Name: "Bob"})
	if result {
		t.Error("Equal should return false for different values")
	}
	if !mt.failed {
		t.Error("Equal should call Errorf for different values")
	}
	if len(mt.messages) == 0 {
		t.Fatal("expected an error message")
	}
	msg := mt.messages[0]
	if !contains(msg, "mismatch (-want +got)") {
		t.Errorf("expected diff header in message, got: %s", msg)
	}
	if !contains(msg, "Alice") || !contains(msg, "Bob") {
		t.Errorf("expected diff to contain values, got: %s", msg)
	}
}

func TestEqual_WithOptions(t *testing.T) {
	mt := &mockT{}
	type S struct {
		Name  string
		Score int
	}
	want := S{Name: "Alice", Score: 100}
	got := S{Name: "Alice", Score: 50}

	// Ignoring Score should make them equal
	opt := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".Score"
	}, cmp.Ignore())

	result := testcmp.Equal(mt, want, got, opt)
	if !result {
		t.Error("Equal should return true when differing field is ignored")
	}
}

func TestEqual_Primitives(t *testing.T) {
	mt := &mockT{}
	if !testcmp.Equal(mt, 42, 42) {
		t.Error("Equal should work with primitive types")
	}

	mt2 := &mockT{}
	if testcmp.Equal(mt2, "hello", "world") {
		t.Error("Equal should return false for different strings")
	}
}

func TestEqual_Slices(t *testing.T) {
	mt := &mockT{}
	want := []int{1, 2, 3}
	got := []int{1, 2, 3}
	if !testcmp.Equal(mt, want, got) {
		t.Error("Equal should return true for identical slices")
	}

	mt2 := &mockT{}
	got2 := []int{1, 2, 4}
	if testcmp.Equal(mt2, want, got2) {
		t.Error("Equal should return false for different slices")
	}
}

func TestNotEqual_DifferentValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.NotEqual(mt, S{Name: "Alice"}, S{Name: "Bob"})
	if !result {
		t.Error("NotEqual should return true for different values")
	}
	if mt.failed {
		t.Error("NotEqual should not call Errorf for different values")
	}
}

func TestNotEqual_SameValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.NotEqual(mt, S{Name: "Alice"}, S{Name: "Alice"})
	if result {
		t.Error("NotEqual should return false for identical values")
	}
	if !mt.failed {
		t.Error("NotEqual should call Errorf for identical values")
	}
	if len(mt.messages) == 0 {
		t.Fatal("expected an error message")
	}
	if !contains(mt.messages[0], "values are equal, but should differ") {
		t.Errorf("unexpected message: %s", mt.messages[0])
	}
}

func TestNotEqual_WithOptions(t *testing.T) {
	mt := &mockT{}
	type S struct {
		Name  string
		Score int
	}
	want := S{Name: "Alice", Score: 100}
	got := S{Name: "Alice", Score: 50}

	// Ignoring Score should make them equal, so NotEqual should fail
	opt := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".Score"
	}, cmp.Ignore())

	result := testcmp.NotEqual(mt, want, got, opt)
	if result {
		t.Error("NotEqual should return false when ignored fields make values equal")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
