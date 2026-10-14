package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/qkhanh39/go-common/apperrors"
)

type TestUser struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"gte=0,lte=150"`
	Password string `json:"password" validate:"required,min=8"`
}

type TestUserWithConfirm struct {
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type TestOptions struct {
	Status string `json:"status" validate:"required,oneof=active inactive pending"`
}

func TestNew(t *testing.T) {
	v := New()
	if v == nil {
		t.Fatal("expected validator instance, got nil")
	}
	if v.validate == nil {
		t.Fatal("expected underlying validator, got nil")
	}
}

func TestValidate(t *testing.T) {
	v := New()

	t.Run("valid struct", func(t *testing.T) {
		user := TestUser{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      30,
			Password: "securepassword123",
		}
		err := v.Validate(user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		user := TestUser{}
		err := v.Validate(user)
		if err == nil {
			t.Error("expected error, got nil")
		}

		valErr, ok := err.(*apperrors.ValidationError)
		if !ok {
			t.Errorf("expected *apperrors.ValidationError, got %T", err)
		}

		// Should have errors for name, email, and password (required fields)
		if len(valErr.Details) < 3 {
			t.Errorf("expected at least 3 errors, got %d", len(valErr.Details))
		}

		// Check that error codes are correct
		fieldCodes := make(map[string]string)
		for _, d := range valErr.Details {
			fieldCodes[d.Field] = d.Code
		}

		if fieldCodes["name"] != "REQUIRED" {
			t.Errorf("expected name code 'REQUIRED', got %q", fieldCodes["name"])
		}
		if fieldCodes["email"] != "REQUIRED" {
			t.Errorf("expected email code 'REQUIRED', got %q", fieldCodes["email"])
		}
		if fieldCodes["password"] != "REQUIRED" {
			t.Errorf("expected password code 'REQUIRED', got %q", fieldCodes["password"])
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		user := TestUser{
			Name:     "John",
			Email:    "not-an-email",
			Password: "securepassword123",
		}
		err := v.Validate(user)
		if err == nil {
			t.Error("expected error, got nil")
		}

		valErr, ok := err.(*apperrors.ValidationError)
		if !ok {
			t.Errorf("expected *apperrors.ValidationError, got %T", err)
		}

		found := false
		for _, d := range valErr.Details {
			if d.Field == "email" && d.Code == "INVALID_EMAIL" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected email field with INVALID_EMAIL code")
		}
	})

	t.Run("min length validation", func(t *testing.T) {
		user := TestUser{
			Name:     "J",
			Email:    "john@example.com",
			Password: "securepassword123",
		}
		err := v.Validate(user)
		if err == nil {
			t.Error("expected error, got nil")
		}

		valErr, ok := err.(*apperrors.ValidationError)
		if !ok {
			t.Errorf("expected *apperrors.ValidationError, got %T", err)
		}

		found := false
		for _, d := range valErr.Details {
			if d.Field == "name" && d.Code == "MIN_LENGTH" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected name field with MIN_LENGTH code")
		}
	})

	t.Run("eqfield validation", func(t *testing.T) {
		user := TestUserWithConfirm{
			Password:        "securepassword123",
			ConfirmPassword: "differentpassword",
		}
		err := v.Validate(user)
		if err == nil {
			t.Error("expected error, got nil")
		}

		valErr, ok := err.(*apperrors.ValidationError)
		if !ok {
			t.Errorf("expected *apperrors.ValidationError, got %T", err)
		}

		found := false
		for _, d := range valErr.Details {
			if d.Field == "confirm_password" && d.Code == "MISMATCH" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected confirm_password field with MISMATCH code")
		}
	})

	t.Run("oneof validation", func(t *testing.T) {
		opts := TestOptions{Status: "invalid"}
		err := v.Validate(opts)
		if err == nil {
			t.Error("expected error, got nil")
		}

		valErr, ok := err.(*apperrors.ValidationError)
		if !ok {
			t.Errorf("expected *apperrors.ValidationError, got %T", err)
		}

		found := false
		for _, d := range valErr.Details {
			if d.Field == "status" && d.Code == "INVALID_OPTION" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected status field with INVALID_OPTION code")
		}
	})

	t.Run("valid oneof", func(t *testing.T) {
		opts := TestOptions{Status: "active"}
		err := v.Validate(opts)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestUnderlying(t *testing.T) {
	v := New()
	underlying := v.Underlying()
	if underlying == nil {
		t.Error("expected underlying validator instance, got nil")
	}
}

func TestRegisterValidation(t *testing.T) {
	v := New()

	// Register a custom validation
	err := v.RegisterValidation("is_even", func(fl validator.FieldLevel) bool {
		return fl.Field().Int()%2 == 0
	})
	if err != nil {
		t.Errorf("expected no error registering validation, got %v", err)
	}
}

func TestDefaultTagCode(t *testing.T) {
	tests := []struct {
		tag      string
		param    string
		expected string
	}{
		{"required", "", "REQUIRED"},
		{"email", "", "INVALID_EMAIL"},
		{"min", "5", "MIN_LENGTH"},
		{"max", "100", "MAX_LENGTH"},
		{"len", "10", "EXACT_LENGTH"},
		{"uuid", "", "INVALID_UUID"},
		{"uuid4", "", "INVALID_UUID"},
		{"url", "", "INVALID_URL"},
		{"eqfield", "Password", "MISMATCH"},
		{"oneof", "a b c", "INVALID_OPTION"},
		{"gt", "0", "GREATER_THAN"},
		{"gte", "0", "GREATER_THAN_OR_EQUAL"},
		{"lt", "100", "LESS_THAN"},
		{"lte", "100", "LESS_THAN_OR_EQUAL"},
		{"numeric", "", "MUST_BE_NUMERIC"},
		{"alpha", "", "MUST_BE_ALPHA"},
		{"alphanum", "", "MUST_BE_ALPHANUMERIC"},
		{"custom_tag", "", "CUSTOM_TAG"},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			result := defaultTagCode(tt.tag, tt.param)
			if result != tt.expected {
				t.Errorf("defaultTagCode(%q, %q) = %q, want %q", tt.tag, tt.param, result, tt.expected)
			}
		})
	}
}

func TestRegisterTagCode(t *testing.T) {
	v := New()

	// Register custom code
	v.RegisterTagCode("required", "FIELD_REQUIRED")

	// Validate struct with required field missing
	type TestStruct struct {
		Name string `json:"name" validate:"required"`
	}

	err := v.Validate(TestStruct{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	valErr, ok := err.(*apperrors.ValidationError)
	if !ok {
		t.Fatalf("expected *apperrors.ValidationError, got %T", err)
	}

	if len(valErr.Details) != 1 {
		t.Fatalf("expected 1 error, got %d", len(valErr.Details))
	}

	if valErr.Details[0].Code != "FIELD_REQUIRED" {
		t.Errorf("expected code 'FIELD_REQUIRED', got %q", valErr.Details[0].Code)
	}
}

func TestRegisterTagMessage(t *testing.T) {
	v := New()

	// Register custom message
	v.RegisterTagMessage("required", func(fe FieldError) string {
		return "this field cannot be empty"
	})

	// Validate struct with required field missing
	type TestStruct struct {
		Name string `json:"name" validate:"required"`
	}

	err := v.Validate(TestStruct{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	valErr, ok := err.(*apperrors.ValidationError)
	if !ok {
		t.Fatalf("expected *apperrors.ValidationError, got %T", err)
	}

	if len(valErr.Details) != 1 {
		t.Fatalf("expected 1 error, got %d", len(valErr.Details))
	}

	if valErr.Details[0].Issue != "this field cannot be empty" {
		t.Errorf("expected issue 'this field cannot be empty', got %q", valErr.Details[0].Issue)
	}
}

func TestChainedRegistration(t *testing.T) {
	v := New().
		RegisterTagCode("required", "MANDATORY").
		RegisterTagMessage("required", func(fe FieldError) string {
			return "this is mandatory"
		})

	type TestStruct struct {
		Name string `json:"name" validate:"required"`
	}

	err := v.Validate(TestStruct{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	valErr := err.(*apperrors.ValidationError)
	if valErr.Details[0].Code != "MANDATORY" {
		t.Errorf("expected code 'MANDATORY', got %q", valErr.Details[0].Code)
	}
	if valErr.Details[0].Issue != "this is mandatory" {
		t.Errorf("expected issue 'this is mandatory', got %q", valErr.Details[0].Issue)
	}
}
