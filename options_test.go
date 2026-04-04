package testcmp_test

import (
	"testing"
	"time"

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
	VALID      string // should not be ignored (doesn't end in "ID")
}

func TestIgnoreTimestamps(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{Name: "Alice", CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	got := sampleStruct{Name: "Alice", CreatedAt: time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	if !result {
		t.Error("IgnoreTimestamps should make structs with different timestamps equal")
	}
}

func TestIgnoreTimestamps_StillComparesOtherFields(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{Name: "Bob"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	if result {
		t.Error("IgnoreTimestamps should not ignore non-timestamp fields")
	}
}

func TestIgnoreAuditFields(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{
		Name: "Alice",
	}
	got := sampleStruct{
		Name:      "Alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "system",
		UpdatedBy: "admin",
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	if !result {
		t.Errorf("IgnoreAuditFields should ignore CreatedAt, UpdatedAt, CreatedBy, UpdatedBy; messages: %v", mt.messages)
	}
}

func TestIgnoreAuditFields_StillComparesOtherFields(t *testing.T) {
	mt := &mockT{}
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{Name: "Bob", CreatedBy: "system"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	if result {
		t.Error("IgnoreAuditFields should not ignore Name")
	}
}

func TestIgnoreIDs(t *testing.T) {
	mt := &mockT{}
	want := withIDs{
		Name:  "Alice",
		VALID: "yes",
	}
	got := withIDs{
		ID:         42,
		SupplierID: "sup-123",
		UserID:     99,
		Name:       "Alice",
		VALID:      "yes",
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	if !result {
		t.Errorf("IgnoreIDs should ignore fields ending in ID; messages: %v", mt.messages)
	}
}

func TestIgnoreIDs_DoesNotIgnoreNonIDFields(t *testing.T) {
	mt := &mockT{}
	want := withIDs{Name: "Alice", VALID: "yes"}
	got := withIDs{Name: "Bob", VALID: "yes"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	if result {
		t.Error("IgnoreIDs should not ignore Name")
	}
}

func TestIgnoreIDs_IgnoresFieldsEndingInID(t *testing.T) {
	// VALID ends in "ID" so it is matched by IgnoreIDs.
	// This is a known trade-off documented in the IgnoreIDs godoc:
	// use IgnoreFields for precise control when convention-based filtering is too broad.
	mt := &mockT{}
	want := withIDs{Name: "Alice", VALID: "yes"}
	got := withIDs{Name: "Alice", VALID: "no"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	if !result {
		t.Error("IgnoreIDs matches any field ending in ID, including VALID")
	}
}

func TestIgnoreFields(t *testing.T) {
	mt := &mockT{}
	type Profile struct {
		Name       string
		ExternalID string
		Score      int
	}
	want := Profile{Name: "Alice", Score: 92}
	got := Profile{Name: "Alice", ExternalID: "ext-123", Score: 92}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreFields(Profile{}, "ExternalID"))
	if !result {
		t.Error("IgnoreFields should ignore the specified fields")
	}
}

func TestIgnoreFields_StillComparesOtherFields(t *testing.T) {
	mt := &mockT{}
	type Profile struct {
		Name       string
		ExternalID string
		Score      int
	}
	want := Profile{Name: "Alice", Score: 92}
	got := Profile{Name: "Alice", ExternalID: "ext-123", Score: 50}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreFields(Profile{}, "ExternalID"))
	if result {
		t.Error("IgnoreFields should still compare non-ignored fields")
	}
}

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
	if !result {
		t.Errorf("Combined options should work together; messages: %v", mt.messages)
	}
}
