package api

import "github.com/qkhanh39/go-common/apperrors"

type ErrorDetail = apperrors.ErrorDetail

type errorInfo struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error errorInfo `json:"error"`
}
