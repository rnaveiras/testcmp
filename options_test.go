package testcmp_test

import (
	"testing"
	"time"

	"github.com/guregu/null/v6"
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

// --- IgnoreAuditFields with pointer time fields ---

func TestIgnoreAuditFields_PointerTimeFields(t *testing.T) {
	// IgnoreAuditFields filters by field name, so it should ignore
	// CreatedAt and UpdatedAt even when they are *time.Time pointers.
	mt := &mockT{}
	now := time.Now()
	want := withPointerTime{Name: "Alice"}
	got := withPointerTime{Name: "Alice", CreatedAt: &now, UpdatedAt: &now}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.True(t, result, "IgnoreAuditFields should ignore *time.Time audit fields by name; messages: %v", mt.messages)
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

// --- Complex nested structs with guregu/null ---

type Address struct {
	ID      int64
	Street  string
	City    string
	ZipCode null.String
}

type Company struct {
	ID        int64
	Name      string
	TaxID     null.String
	Active    null.Bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Employee struct {
	ID         int64
	EmployeeID string
	Name       string
	Email      null.String
	Age        null.Int
	Score      null.Float
	Address    Address
	Company    Company
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CreatedBy  string
	UpdatedBy  string
}

func TestComplexNestedStruct_WithNullFields(t *testing.T) {
	mt := &mockT{}
	want := Employee{
		Name:  "Alice",
		Email: null.StringFrom("alice@example.com"),
		Age:   null.IntFrom(30),
		Score: null.FloatFrom(95.5),
		Address: Address{
			Street:  "123 Main St",
			City:    "London",
			ZipCode: null.StringFrom("EC1A 1BB"),
		},
		Company: Company{
			Name:   "Acme Corp",
			TaxID:  null.StringFrom("GB123456789"),
			Active: null.BoolFrom(true),
		},
	}
	got := Employee{
		ID:         1,
		EmployeeID: "emp-001",
		Name:       "Alice",
		Email:      null.StringFrom("alice@example.com"),
		Age:        null.IntFrom(30),
		Score:      null.FloatFrom(95.5),
		Address: Address{
			ID:      10,
			Street:  "123 Main St",
			City:    "London",
			ZipCode: null.StringFrom("EC1A 1BB"),
		},
		Company: Company{
			ID:        20,
			Name:      "Acme Corp",
			TaxID:     null.StringFrom("GB123456789"),
			Active:    null.BoolFrom(true),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "system",
		UpdatedBy: "admin",
	}

	result := testcmp.Equal(mt, want, got,
		testcmp.IgnoreAuditFields(),
		testcmp.IgnoreIDs(),
	)
	assert.True(t, result, "should ignore IDs and audit fields in complex nested structs with null fields; messages: %v", mt.messages)
}

func TestComplexNestedStruct_NullVsZeroValue(t *testing.T) {
	mt := &mockT{}
	want := Employee{
		Name:  "Alice",
		Email: null.StringFrom("alice@example.com"),
		Age:   null.IntFrom(30),
	}
	got := Employee{
		Name:  "Alice",
		Email: null.StringFrom("alice@example.com"),
		Age:   null.Int{}, // null/zero value — should differ from 30
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs(), testcmp.IgnoreAuditFields())
	assert.False(t, result, "should detect difference between null.IntFrom(30) and zero null.Int")
}

func TestComplexNestedStruct_DifferentNullStringValues(t *testing.T) {
	mt := &mockT{}
	want := Employee{
		Name:    "Alice",
		Address: Address{City: "London", ZipCode: null.StringFrom("EC1A 1BB")},
	}
	got := Employee{
		Name:    "Alice",
		Address: Address{City: "London", ZipCode: null.StringFrom("SW1A 1AA")},
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs(), testcmp.IgnoreAuditFields())
	assert.False(t, result, "should detect different null.String values in nested structs")
}

func TestComplexNestedStruct_OnlyNonAuditFieldsDiffer(t *testing.T) {
	mt := &mockT{}
	want := Employee{
		Name:  "Alice",
		Email: null.StringFrom("alice@example.com"),
		Company: Company{
			Name: "Acme Corp",
		},
	}
	got := Employee{
		Name:  "Alice",
		Email: null.StringFrom("alice@example.com"),
		Company: Company{
			Name: "Globex",
		},
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs(), testcmp.IgnoreAuditFields())
	assert.False(t, result, "should detect company name difference in nested struct")
}
