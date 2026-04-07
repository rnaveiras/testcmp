package testcmp_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rnaveiras/testcmp"
)

func newMockT(t *testing.T) *MockTestingT {
	t.Helper()
	mt := NewMockTestingT(t)
	mt.EXPECT().Helper().Return().Maybe()
	return mt
}

// --- Equal ---

func TestEqual_SameValues(t *testing.T) {
	mt := newMockT(t)
	type S struct{ Name string }
	result := testcmp.Equal(mt, S{Name: "Alice"}, S{Name: "Alice"})
	assert.True(t, result)
}

func TestEqual_DifferentValues(t *testing.T) {
	mt := newMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	type S struct{ Name string }
	result := testcmp.Equal(mt, S{Name: "Alice"}, S{Name: "Bob"})
	assert.False(t, result)
}

func TestEqual_DifferentValues_DiffContent(t *testing.T) {
	var captured string
	mt := newMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).
		Run(func(format string, args ...interface{}) {
			captured = fmt.Sprintf(format, args...)
		}).Return().Once()

	type S struct{ Name string }
	testcmp.Equal(mt, S{Name: "Alice"}, S{Name: "Bob"})

	assert.Contains(t, captured, "mismatch (-want +got)")
}

func TestEqual_WithOptions(t *testing.T) {
	mt := newMockT(t)
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
	assert.True(t, result)
}

func TestEqual_Primitives(t *testing.T) {
	mt := newMockT(t)
	assert.True(t, testcmp.Equal(mt, 42, 42))

	mt2 := newMockT(t)
	mt2.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	assert.False(t, testcmp.Equal(mt2, "hello", "world"))
}

func TestEqual_Slices(t *testing.T) {
	mt := newMockT(t)
	assert.True(t, testcmp.Equal(mt, []int{1, 2, 3}, []int{1, 2, 3}))

	mt2 := newMockT(t)
	mt2.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	assert.False(t, testcmp.Equal(mt2, []int{1, 2, 3}, []int{1, 2, 4}))
}

func TestEqual_NilVsEmptySlice(t *testing.T) {
	mt := newMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	var nilSlice []int
	emptySlice := []int{}

	result := testcmp.Equal(mt, nilSlice, emptySlice)
	assert.False(t, result)
}

func TestEqual_DiffWithPercentChar(t *testing.T) {
	var captured string
	mt := newMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).
		Run(func(format string, args ...interface{}) {
			captured = fmt.Sprintf(format, args...)
		}).Return().Once()

	testcmp.Equal(mt, "100%", "50%")
	assert.Contains(t, captured, "100%")
	assert.Contains(t, captured, "50%")
}

// --- NotEqual ---

func TestNotEqual_DifferentValues(t *testing.T) {
	mt := newMockT(t)
	type S struct{ Name string }
	result := testcmp.NotEqual(mt, S{Name: "Alice"}, S{Name: "Bob"})
	assert.True(t, result)
}

func TestNotEqual_SameValues(t *testing.T) {
	mt := newMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string")).Return().Once()
	type S struct{ Name string }
	result := testcmp.NotEqual(mt, S{Name: "Alice"}, S{Name: "Alice"})
	assert.False(t, result)
}

func TestNotEqual_WithOptions(t *testing.T) {
	mt := newMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string")).Return().Once()
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
	assert.False(t, result)
}

// --- Maps ---

func TestEqual_Maps(t *testing.T) {
	mt := newMockT(t)
	want := map[string]int{"a": 1, "b": 2}
	got := map[string]int{"a": 1, "b": 2}
	assert.True(t, testcmp.Equal(mt, want, got))

	mt2 := newMockT(t)
	mt2.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	got2 := map[string]int{"a": 1, "b": 3}
	assert.False(t, testcmp.Equal(mt2, want, got2))
}

func TestNotEqual_Maps(t *testing.T) {
	mt := newMockT(t)
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

func TestEqual_WithTBRecorder(t *testing.T) {
	rec := newTBRecorder(t)
	type S struct{ Name string }
	result := testcmp.Equal(rec, S{Name: "Alice"}, S{Name: "Alice"})
	assert.True(t, result)
	assert.False(t, rec.failed)
}
