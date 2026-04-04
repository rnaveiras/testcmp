package testcmp_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rnaveiras/testcmp"
)

// --- Equal ---

func TestEqual_SameValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.Equal(mt, S{Name: "Alice"}, S{Name: "Alice"})
	assert.True(t, result, "Equal should return true for identical values")
	assert.False(t, mt.failed, "Equal should not call Errorf for identical values")
}

func TestEqual_DifferentValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.Equal(mt, S{Name: "Alice"}, S{Name: "Bob"})
	assert.False(t, result, "Equal should return false for different values")
	assert.True(t, mt.failed, "Equal should call Errorf for different values")
	require.NotEmpty(t, mt.messages)
	assert.Contains(t, mt.messages[0], "mismatch (-want +got)")
	assert.Contains(t, mt.messages[0], "Alice")
	assert.Contains(t, mt.messages[0], "Bob")
}

func TestEqual_WithOptions(t *testing.T) {
	mt := &mockT{}
	type S struct {
		Name  string
		Score int
	}
	want := S{Name: "Alice", Score: 100}
	got := S{Name: "Alice", Score: 50}

	opt := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".Score"
	}, cmp.Ignore())

	result := testcmp.Equal(mt, want, got, opt)
	assert.True(t, result, "Equal should return true when differing field is ignored")
}

func TestEqual_Primitives(t *testing.T) {
	mt := &mockT{}
	assert.True(t, testcmp.Equal(mt, 42, 42))

	mt2 := &mockT{}
	assert.False(t, testcmp.Equal(mt2, "hello", "world"))
}

func TestEqual_Slices(t *testing.T) {
	mt := &mockT{}
	assert.True(t, testcmp.Equal(mt, []int{1, 2, 3}, []int{1, 2, 3}))

	mt2 := &mockT{}
	assert.False(t, testcmp.Equal(mt2, []int{1, 2, 3}, []int{1, 2, 4}))
}

func TestEqual_NilVsEmptySlice(t *testing.T) {
	mt := &mockT{}
	var nilSlice []int
	emptySlice := []int{}

	result := testcmp.Equal(mt, nilSlice, emptySlice)
	assert.False(t, result, "Equal should distinguish nil slice from empty slice")
}

func TestEqual_DiffWithPercentChar(t *testing.T) {
	mt := &mockT{}
	testcmp.Equal(mt, "100%", "50%")
	require.True(t, mt.failed)
	assert.Contains(t, mt.messages[0], "100%")
	assert.Contains(t, mt.messages[0], "50%")
}

// --- NotEqual ---

func TestNotEqual_DifferentValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.NotEqual(mt, S{Name: "Alice"}, S{Name: "Bob"})
	assert.True(t, result, "NotEqual should return true for different values")
	assert.False(t, mt.failed, "NotEqual should not call Errorf for different values")
}

func TestNotEqual_SameValues(t *testing.T) {
	mt := &mockT{}
	type S struct{ Name string }
	result := testcmp.NotEqual(mt, S{Name: "Alice"}, S{Name: "Alice"})
	assert.False(t, result, "NotEqual should return false for identical values")
	assert.True(t, mt.failed, "NotEqual should call Errorf for identical values")
	require.NotEmpty(t, mt.messages)
	assert.Contains(t, mt.messages[0], "values are equal, but should differ")
}

func TestNotEqual_WithOptions(t *testing.T) {
	mt := &mockT{}
	type S struct {
		Name  string
		Score int
	}
	want := S{Name: "Alice", Score: 100}
	got := S{Name: "Alice", Score: 50}

	opt := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".Score"
	}, cmp.Ignore())

	result := testcmp.NotEqual(mt, want, got, opt)
	assert.False(t, result, "NotEqual should return false when ignored fields make values equal")
}

// --- Maps ---

func TestEqual_Maps(t *testing.T) {
	mt := &mockT{}
	want := map[string]int{"a": 1, "b": 2}
	got := map[string]int{"a": 1, "b": 2}
	assert.True(t, testcmp.Equal(mt, want, got))

	mt2 := &mockT{}
	got2 := map[string]int{"a": 1, "b": 3}
	assert.False(t, testcmp.Equal(mt2, want, got2))
}

func TestNotEqual_Maps(t *testing.T) {
	mt := &mockT{}
	want := map[string]int{"a": 1}
	got := map[string]int{"a": 2}
	assert.True(t, testcmp.NotEqual(mt, want, got))
}

// --- Real *testing.T integration ---

func TestEqual_WithRealTestingT(t *testing.T) {
	type S struct{ Name string }
	testcmp.Equal(t, S{Name: "Alice"}, S{Name: "Alice"})
}

func TestNotEqual_WithRealTestingT(t *testing.T) {
	type S struct{ Name string }
	testcmp.NotEqual(t, S{Name: "Alice"}, S{Name: "Bob"})
}
