package apperrors

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	CodeValidationFailed = "VALIDATION_FAILED"
	CodeBadRequest       = "BAD_REQUEST"
	CodeUnauthorized     = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeNotFound         = "NOT_FOUND"
	CodeConflict         = "CONFLICT"
	CodeInternalError    = "INTERNAL_ERROR"
)

type AppError interface {
	error
	Status() int
	ErrCode() string
}

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Code  string `json:"code,omitempty"`
	Issue string `json:"issue"`
}

type ValidationError struct {
	Details []ErrorDetail `json:"details"`
}

func (e *ValidationError) Error() string {
	if len(e.Details) == 0 {
		return "validation failed"
	}

	var msgs []string
	for _, d := range e.Details {
		if d.Field != "" {
			msgs = append(msgs, fmt.Sprintf("%s: %s", d.Field, d.Issue))
		} else {
			msgs = append(msgs, d.Issue)
		}
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, "; "))
}

func (e *ValidationError) Status() int {
	return http.StatusUnprocessableEntity
}

func (e *ValidationError) ErrCode() string {
	return CodeValidationFailed
}

func (e *ValidationError) Add(field, issue, code string) {
	e.Details = append(e.Details, ErrorDetail{Field: field, Issue: issue, Code: code})
}

func (e *ValidationError) AddGeneral(issue, code string) {
	e.Details = append(e.Details, ErrorDetail{Issue: issue, Code: code})
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Details) > 0
}

func NewValidationError() *ValidationError {
	return &ValidationError{}
}

type BadRequestError ErrorDetail

func (e BadRequestError) Error() string   { return e.Issue }
func (e BadRequestError) Status() int     { return http.StatusBadRequest }
func (e BadRequestError) ErrCode() string { return CodeBadRequest }

type UnauthorizedError ErrorDetail

func (e UnauthorizedError) Error() string   { return e.Issue }
func (e UnauthorizedError) Status() int     { return http.StatusUnauthorized }
func (e UnauthorizedError) ErrCode() string { return CodeUnauthorized }

type ForbiddenError ErrorDetail

func (e ForbiddenError) Error() string   { return e.Issue }
func (e ForbiddenError) Status() int     { return http.StatusForbidden }
func (e ForbiddenError) ErrCode() string { return CodeForbidden }

type NotFoundError ErrorDetail

func (e NotFoundError) Error() string   { return e.Issue }
func (e NotFoundError) Status() int     { return http.StatusNotFound }
func (e NotFoundError) ErrCode() string { return CodeNotFound }

type ConflictError ErrorDetail

func (e ConflictError) Error() string   { return e.Issue }
func (e ConflictError) Status() int     { return http.StatusConflict }
func (e ConflictError) ErrCode() string { return CodeConflict }

type InternalError ErrorDetail

func (e InternalError) Error() string   { return e.Issue }
func (e InternalError) Status() int     { return http.StatusInternalServerError }
func (e InternalError) ErrCode() string { return CodeInternalError }
