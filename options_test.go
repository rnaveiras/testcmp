package testcmp_test

import (
	"testing"
	"time"

	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rnaveiras/testcmp"
	"github.com/rnaveiras/testcmp/mocks"
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

func newOptionMockT(t *testing.T) *mocks.MockTestingT {
	t.Helper()
	mt := mocks.NewMockTestingT(t)
	mt.EXPECT().Helper().Return().Maybe()
	return mt
}

// --- IgnoreTimestamps ---

func TestIgnoreTimestamps(t *testing.T) {
	mt := newOptionMockT(t)
	want := sampleStruct{Name: "Alice", CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	got := sampleStruct{Name: "Alice", CreatedAt: time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	assert.True(t, result)
}

func TestIgnoreTimestamps_StillComparesOtherFields(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{Name: "Bob"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	assert.False(t, result)
}

// --- IgnoreAuditFields ---

func TestIgnoreAuditFields(t *testing.T) {
	mt := newOptionMockT(t)
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{
		Name:      "Alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "system",
		UpdatedBy: "admin",
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.True(t, result)
}

func TestIgnoreAuditFields_StillComparesOtherFields(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	want := sampleStruct{Name: "Alice"}
	got := sampleStruct{Name: "Bob", CreatedBy: "system"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.False(t, result)
}

func TestIgnoreAuditFields_PreservesNonAuditTimeFields(t *testing.T) {
	type event struct {
		Name        string
		ScheduledAt time.Time
		CreatedAt   time.Time
	}
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	want := event{Name: "deploy", ScheduledAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	got := event{
		Name:        "deploy",
		ScheduledAt: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.False(t, result)
}

func TestIgnoreAuditFields_Nested(t *testing.T) {
	mt := newOptionMockT(t)
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
	assert.True(t, result)
}

// --- IgnoreIDs ---

func TestIgnoreIDs(t *testing.T) {
	mt := newOptionMockT(t)
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
	assert.True(t, result)
}

func TestIgnoreIDs_DoesNotIgnoreNonIDFields(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	want := withIDs{Name: "Alice"}
	got := withIDs{Name: "Bob"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.False(t, result)
}

func TestIgnoreIDs_DoesNotMatchVALID(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	want := withIDs{Name: "Alice", VALID: "yes"}
	got := withIDs{Name: "Alice", VALID: "no"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.False(t, result)
}

func TestIgnoreIDs_DoesNotMatchSQUID(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	want := withIDs{Name: "Alice", SQUID: "yes"}
	got := withIDs{Name: "Alice", SQUID: "no"}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.False(t, result)
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
	mt := newOptionMockT(t)
	want := outer{Inner: inner{Name: "x"}}
	got := outer{ID: 1, Inner: inner{ID: 2, Name: "x"}}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs())
	assert.True(t, result)
}

// --- IgnoreAuditFields with pointer time fields ---

func TestIgnoreAuditFields_PointerTimeFields(t *testing.T) {
	mt := newOptionMockT(t)
	now := time.Now()
	want := withPointerTime{Name: "Alice"}
	got := withPointerTime{Name: "Alice", CreatedAt: &now, UpdatedAt: &now}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreAuditFields())
	assert.True(t, result)
}

// --- IgnoreTimestamps with pointer fields ---

func TestIgnoreTimestamps_DoesNotIgnorePointerTimeFields(t *testing.T) {
	// IgnoreTimestamps uses cmpopts.IgnoreTypes(time.Time{}) which matches
	// time.Time but not *time.Time. This is expected — use cmpopts.IgnoreFields
	// for pointer time fields if needed.
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	now := time.Now()
	want := withPointerTime{Name: "Alice"}
	got := withPointerTime{Name: "Alice", CreatedAt: &now, UpdatedAt: &now}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreTimestamps())
	assert.False(t, result)
}

// --- Combined options ---

func TestCombinedOptions(t *testing.T) {
	mt := newOptionMockT(t)
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
	assert.True(t, result)
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
	mt := newOptionMockT(t)
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
	assert.True(t, result)
}

func TestComplexNestedStruct_NullVsZeroValue(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
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
	assert.False(t, result)
}

func TestComplexNestedStruct_DifferentNullStringValues(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
	want := Employee{
		Name:    "Alice",
		Address: Address{City: "London", ZipCode: null.StringFrom("EC1A 1BB")},
	}
	got := Employee{
		Name:    "Alice",
		Address: Address{City: "London", ZipCode: null.StringFrom("SW1A 1AA")},
	}

	result := testcmp.Equal(mt, want, got, testcmp.IgnoreIDs(), testcmp.IgnoreAuditFields())
	assert.False(t, result)
}

func TestComplexNestedStruct_OnlyNonAuditFieldsDiffer(t *testing.T) {
	mt := newOptionMockT(t)
	mt.EXPECT().Errorf(mock.AnythingOfType("string"), mock.Anything).Return().Once()
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
	assert.False(t, result)
}
