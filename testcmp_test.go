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
	rec := newTBRecorder(t)
	type S struct{ Name string }
	result := testcmp.Equal(rec, S{Name: "Alice"}, S{Name: "Alice"})
	assert.True(t, result)
	assert.False(t, rec.failed)
}

func TestEqual_DifferentValues(t *testing.T) {
	rec := newTBRecorder(t)
	type S struct{ Name string }
	result := testcmp.Equal(rec, S{Name: "Alice"}, S{Name: "Bob"})
	assert.False(t, result)
	assert.True(t, rec.failed)
	require.Len(t, rec.errors, 1)
}

func TestEqual_DifferentValues_DiffContent(t *testing.T) {
	rec := newTBRecorder(t)
	type S struct{ Name string }
	testcmp.Equal(rec, S{Name: "Alice"}, S{Name: "Bob"})

	require.Len(t, rec.errors, 1)
	assert.Contains(t, rec.errors[0], "mismatch (-want +got)")
}

func TestEqual_WithOptions(t *testing.T) {
	rec := newTBRecorder(t)
	type S struct {
		Name  string
		Score int
	}
	want := S{Name: "Alice", Score: 100}
	got := S{Name: "Alice", Score: 50}

	opt := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".Score"
	}, cmp.Ignore())

	result := testcmp.Equal(rec, want, got, opt)
	assert.True(t, result)
	assert.False(t, rec.failed)
}

func TestEqual_Primitives(t *testing.T) {
	rec := newTBRecorder(t)
	assert.True(t, testcmp.Equal(rec, 42, 42))
	assert.False(t, rec.failed)

	rec2 := newTBRecorder(t)
	assert.False(t, testcmp.Equal(rec2, "hello", "world"))
	assert.True(t, rec2.failed)
}

func TestEqual_Slices(t *testing.T) {
	rec := newTBRecorder(t)
	assert.True(t, testcmp.Equal(rec, []int{1, 2, 3}, []int{1, 2, 3}))
	assert.False(t, rec.failed)

	rec2 := newTBRecorder(t)
	assert.False(t, testcmp.Equal(rec2, []int{1, 2, 3}, []int{1, 2, 4}))
	assert.True(t, rec2.failed)
}

func TestEqual_NilVsEmptySlice(t *testing.T) {
	rec := newTBRecorder(t)
	var nilSlice []int
	emptySlice := []int{}

	result := testcmp.Equal(rec, nilSlice, emptySlice)
	assert.False(t, result)
	assert.True(t, rec.failed)
}

func TestEqual_DiffWithPercentChar(t *testing.T) {
	rec := newTBRecorder(t)
	testcmp.Equal(rec, "100%", "50%")

	require.Len(t, rec.errors, 1)
	assert.Contains(t, rec.errors[0], "100%")
	assert.Contains(t, rec.errors[0], "50%")
}

// --- NotEqual ---

func TestNotEqual_DifferentValues(t *testing.T) {
	rec := newTBRecorder(t)
	type S struct{ Name string }
	result := testcmp.NotEqual(rec, S{Name: "Alice"}, S{Name: "Bob"})
	assert.True(t, result)
	assert.False(t, rec.failed)
}

func TestNotEqual_SameValues(t *testing.T) {
	rec := newTBRecorder(t)
	type S struct{ Name string }
	result := testcmp.NotEqual(rec, S{Name: "Alice"}, S{Name: "Alice"})
	assert.False(t, result)
	assert.True(t, rec.failed)
}

func TestNotEqual_WithOptions(t *testing.T) {
	rec := newTBRecorder(t)
	type S struct {
		Name  string
		Score int
	}
	want := S{Name: "Alice", Score: 100}
	got := S{Name: "Alice", Score: 50}

	opt := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".Score"
	}, cmp.Ignore())

	result := testcmp.NotEqual(rec, want, got, opt)
	assert.False(t, result)
	assert.True(t, rec.failed)
}

// --- Maps ---

func TestEqual_Maps(t *testing.T) {
	rec := newTBRecorder(t)
	want := map[string]int{"a": 1, "b": 2}
	got := map[string]int{"a": 1, "b": 2}
	assert.True(t, testcmp.Equal(rec, want, got))
	assert.False(t, rec.failed)

	rec2 := newTBRecorder(t)
	got2 := map[string]int{"a": 1, "b": 3}
	assert.False(t, testcmp.Equal(rec2, want, got2))
	assert.True(t, rec2.failed)
}

func TestNotEqual_Maps(t *testing.T) {
	rec := newTBRecorder(t)
	want := map[string]int{"a": 1}
	got := map[string]int{"a": 2}
	assert.True(t, testcmp.NotEqual(rec, want, got))
	assert.False(t, rec.failed)
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
