package testcmp_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rnaveiras/testcmp"
)

type sampleStruct struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
}

type withIDs struct {
	ID         int64
	SupplierID string
	UserID     int64
	Name       string
	VALID      string
	SQUID      string
}

type nested struct {
	Name  string
	Inner sampleStruct
}

type withPointerTime struct {
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// --- IgnoreTimestamps ---

func TestIgnoreTimestamps(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{Name: "Alice", CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	got := sampleStruct{Name: "Alice", CreatedAt: time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	assert.True(t, result, "IgnoreTimestamps should make structs with different timestamps equal")
}

func TestIgnoreTimestamps_StillComparesOtherFields(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{Name: "Bob"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	assert.False(t, result, "IgnoreTimestamps should not ignore non-timestamp fields")
}

// --- IgnoreAuditFields ---

func TestIgnoreAuditFields(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{
		Name:      "Alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "system",
		UpdatedBy: "admin",
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.True(t, result, "IgnoreAuditFields should ignore CreatedAt, UpdatedAt, CreatedBy, UpdatedBy; messages: %v", mt.messages)
}

func TestIgnoreAuditFields_StillComparesOtherFields(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{Name: "Bob", CreatedBy: "system"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.False(t, result, "IgnoreAuditFields should not ignore Name")
}

func TestIgnoreAuditFields_PreservesNonAuditTimeFields(t *testing.T) {
	type event struct {
		Name        string
		ScheduledAt time.Time
		CreatedAt   time.Time
	}
	mt := &mockT{}
	want := event{Name: "deploy", ScheduledAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	got := event{
		Name:        "deploy",
		ScheduledAt: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.False(t, result, "IgnoreAuditFields should still compare ScheduledAt (non-audit time field)")
}

func TestIgnoreAuditFields_Nested(t *testing.T) {
	mt := &mockT{}
	want := nested{
		Name:  "outer",
		Inner: sampleStruct{Name: "inner"},
	}
	got := nested{
		Name: "outer",
		Inner: sampleStruct{
			Name:      "inner",
			CreatedAt: time.Now(),
			UpdatedBy: "admin",
		},
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.True(t, result, "IgnoreAuditFields should work on nested structs; messages: %v", mt.messages)
}

// --- IgnoreIDs ---

func TestIgnoreIDs(t *testing.T) {
	mt := &mockT{}
	want := withIDs{Name: "Alice", VALID: "yes", SQUID: "yes"}
	got := withIDs{
		ID:         42,
		SupplierID: "sup-123",
		UserID:     99,
		Name:       "Alice",
		VALID:      "yes",
		SQUID:      "yes",
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.True(t, result, "IgnoreIDs should ignore ID, SupplierID, UserID; messages: %v", mt.messages)
}

func TestIgnoreIDs_DoesNotIgnoreNonIDFields(t *testing.T) {
	mt := &mockT{}
	want := withIDs{Name: "Alice"}
	got := withIDs{Name: "Bob"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.False(t, result, "IgnoreIDs should not ignore Name")
}

func TestIgnoreIDs_DoesNotMatchVALID(t *testing.T) {
	mt := &mockT{}
	want := withIDs{Name: "Alice", VALID: "yes"}
	got := withIDs{Name: "Alice", VALID: "no"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.False(t, result, "IgnoreIDs should not match VALID (uppercase before ID)")
}

func TestIgnoreIDs_DoesNotMatchSQUID(t *testing.T) {
	mt := &mockT{}
	want := withIDs{Name: "Alice", SQUID: "yes"}
	got := withIDs{Name: "Alice", SQUID: "no"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.False(t, result, "IgnoreIDs should not match SQUID (uppercase before ID)")
}

func TestIgnoreIDs_Nested(t *testing.T) {
	type inner struct {
		ID   int64
		Name string
	}
	type outer struct {
		ID    int64
		Inner inner
	}
	mt := &mockT{}
	want := outer{Inner: inner{Name: "x"}}
	got := outer{ID: 1, Inner: inner{ID: 2, Name: "x"}}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.True(t, result, "IgnoreIDs should work on nested structs; messages: %v", mt.messages)
}

// --- IgnoreTimestamps with pointer fields ---

func TestIgnoreTimestamps_DoesNotIgnorePointerTimeFields(t *testing.T) {
	// IgnoreTimestamps uses cmpopts.IgnoreTypes(time.Time{}) which matches
	// time.Time but not *time.Time. This is expected — use cmpopts.IgnoreFields
	// for pointer time fields if needed.
	mt := &mockT{}
	now := time.Now()
	want := withPointerTime{Name: "Alice"}
	got := withPointerTime{Name: "Alice", CreatedAt: &now, UpdatedAt: &now}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	assert.False(t, result, "IgnoreTimestamps should not ignore *time.Time pointer fields")
}

// --- Combined options ---

func TestCombinedOptions(t *testing.T) {
	mt := &mockT{}
	type Record struct {
		ID        int64
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
		CreatedBy string
		UpdatedBy string
	}
	want := Record{Name: "Alice"}
	got := Record{
		ID:        1,
		Name:      "Alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "system",
		UpdatedBy: "admin",
	}

	result := testcmp.Equal(mt, want, got,
		testcmp.IgnoreAuditFields(),
		testcmp.IgnoreIDs(),
	)
	assert.True(t, result, "Combined options should work together; messages: %v", mt.messages)
}
