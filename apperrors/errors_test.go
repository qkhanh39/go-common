package apperrors

import (
	"net/http"
	"testing"
)

func TestValidationError(t *testing.T) {
	t.Run("empty validation error", func(t *testing.T) {
		err := NewValidationError()
		if err.Error() != "validation failed" {
			t.Errorf("expected 'valdation failed', got %q", err.Error())
		}
		if err.Status() != http.StatusUnprocessableEntity {
			t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, err.Status())
		}
		if err.ErrCode() != "VALIDATION_FAILED" {
			t.Errorf("expected code 'VALIDATION_FAILED', got %q", err.ErrCode())
		}
		if err.HasErrors() {
			t.Error("expected HasError() to return false")
		}
	})

	t.Run("validation error with field details", func(t *testing.T) {
		err := NewValidationError()
		err.Add("email", "must be a valid email", "INVALID_EMAIL")
		err.Add("name", "is required", "REQUIRED")

		if !err.HasErrors() {
			t.Error("expected HasErrors() to return true")
		}
		if len(err.Details) != 2 {
			t.Errorf("expected 2 details, got %d", len(err.Details))
		}

		expected := "validation failed: email: must be a valid email; name: is required"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("validation error with general detail", func(t *testing.T) {
		err := NewValidationError()
		err.AddGeneral("passwords do not match", "MISMATCH")

		if len(err.Details) != 1 {
			t.Errorf("expected 1 detail, got %d", len(err.Details))
		}
		if err.Details[0].Field != "" {
			t.Errorf("expected empty field, got %q", err.Details[0].Field)
		}

		expected := "validation failed: passwords do not match"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})
}

func TestBadRequestError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := &BadRequestError{Issue: "invalid JSON"}
		if err.Error() != "invalid JSON" {
			t.Errorf("expected 'invalid JSON', got %q", err.Error())
		}
		if err.Status() != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, err.Status())
		}
		if err.ErrCode() != "BAD_REQUEST" {
			t.Errorf("expected code 'BAD_REQUEST', got %q", err.ErrCode())
		}
	})

}

func TestUnauthorizedError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := &UnauthorizedError{Issue: "invalid token"}
		if err.Error() != "invalid token" {
			t.Errorf("expected 'invalid token', got %q", err.Error())
		}
		if err.Status() != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, err.Status())
		}
		if err.ErrCode() != "UNAUTHORIZED" {
			t.Errorf("expected code 'UNAUTHORIZED', got %q", err.ErrCode())
		}
	})

}

func TestForbiddenError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := &ForbiddenError{Issue: "access denied"}
		if err.Error() != "access denied" {
			t.Errorf("expected 'access denied', got %q", err.Error())
		}
		if err.Status() != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, err.Status())
		}
		if err.ErrCode() != "FORBIDDEN" {
			t.Errorf("expected code 'FORBIDDEN', got %q", err.ErrCode())
		}
	})

}

func TestNotFoundError(t *testing.T) {
	t.Run("with resource and id", func(t *testing.T) {
		err := &NotFoundError{Field: "user", Issue: "user not found: 123"}
		if err.Error() != "user not found: 123" {
			t.Errorf("expected 'user not found: 123', got %q", err.Error())
		}
		if err.Status() != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, err.Status())
		}
		if err.ErrCode() != "NOT_FOUND" {
			t.Errorf("expected code 'NOT_FOUND', got %q", err.ErrCode())
		}
	})

	t.Run("with resource only", func(t *testing.T) {
		err := &NotFoundError{Field: "user", Issue: "user not found"}
		if err.Error() != "user not found" {
			t.Errorf("expected 'user not found', got %q", err.Error())
		}
	})
}

func TestConflictError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := &ConflictError{Issue: "email already exists"}
		if err.Error() != "email already exists" {
			t.Errorf("expected 'email already exists', got %q", err.Error())
		}
		if err.Status() != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, err.Status())
		}
		if err.ErrCode() != "CONFLICT" {
			t.Errorf("expected code 'CONFLICT', got %q", err.ErrCode())
		}
	})

}

func TestInternalError(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		err := &InternalError{Issue: "something went wrong"}
		if err.Error() != "something went wrong" {
			t.Errorf("expected 'something went wrong', got %q", err.Error())
		}
		if err.Status() != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, err.Status())
		}
		if err.ErrCode() != "INTERNAL_ERROR" {
			t.Errorf("expected code 'INTERNAL_ERROR', got %q", err.ErrCode())
		}
	})

}

func TestAppErrorInterface(t *testing.T) {
	// Ensure all error types implement AppError interface
	var _ AppError = &ValidationError{}
	var _ AppError = &BadRequestError{}
	var _ AppError = &UnauthorizedError{}
	var _ AppError = &ForbiddenError{}
	var _ AppError = &NotFoundError{}
	var _ AppError = &ConflictError{}
	var _ AppError = &InternalError{}
}
