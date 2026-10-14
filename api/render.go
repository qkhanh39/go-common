package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/qkhanh39/go-common/apperrors"
)

func RenderJSON(w http.ResponseWriter, status int, v any, headers map[string]string) {
	if status <= 0 {
		status = http.StatusOK
	}

	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}

	var buf bytes.Buffer

	if v != nil {
		if err := json.NewEncoder(&buf).Encode(v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(status)

	if v != nil {
		_, _ = w.Write(buf.Bytes())
	}
}

type ResponseBuilder struct {
	w      http.ResponseWriter
	r      *http.Request
	status int

	data     any
	messages []string
	headers  map[string]string

	errCode    string
	errMessage string
	errDetails []ErrorDetail
}

func Response(w http.ResponseWriter, r *http.Request) *ResponseBuilder {
	return &ResponseBuilder{w: w, r: r, status: http.StatusOK}
}

func (b *ResponseBuilder) Status(status int) *ResponseBuilder {
	b.status = status
	return b
}

func (b *ResponseBuilder) WithData(data any) *ResponseBuilder {
	b.data = data
	return b
}

func (b *ResponseBuilder) WithMessages(msgs ...string) *ResponseBuilder {
	b.messages = msgs
	return b
}

func (b *ResponseBuilder) WithHeader(key, value string) *ResponseBuilder {
	if b.headers == nil {
		b.headers = make(map[string]string)
	}
	b.headers[key] = value
	return b
}

func (b *ResponseBuilder) WithHeaders(headers map[string]string) *ResponseBuilder {
	if b.headers == nil {
		b.headers = make(map[string]string)
	}
	for key, value := range headers {
		b.headers[key] = value
	}
	return b
}

func (b *ResponseBuilder) Code(code string) *ResponseBuilder {
	b.errCode = code
	return b
}

func (b *ResponseBuilder) ErrorMessage(message string) *ResponseBuilder {
	b.errMessage = message
	return b
}

func (b *ResponseBuilder) WithErrorDetails(details []ErrorDetail) *ResponseBuilder {
	b.errDetails = details
	return b
}

func (b *ResponseBuilder) WithError(err error) *ResponseBuilder {
	if err == nil {
		return b
	}

	var appErr apperrors.AppError
	if !errors.As(err, &appErr) {
		// Unknown error type - return generic internal error
		b.status = http.StatusInternalServerError
		b.errCode = apperrors.CodeInternalError
		b.errMessage = err.Error()
		return b
	}

	b.status = appErr.Status()
	if b.status <= 0 {
		b.status = http.StatusInternalServerError
	}
	b.errCode = appErr.ErrCode()

	// Special handling for ValidationError — reuse details directly
	if valErr, ok := err.(*apperrors.ValidationError); ok {
		b.errMessage = "Validation failed"
		b.errDetails = valErr.Details
		return b
	}

	b.errMessage = err.Error()
	return b
}

func (b *ResponseBuilder) Send() {
	if b.errCode != "" {
		b.sendError()
	} else {
		b.sendSuccess()
	}
}

func (b *ResponseBuilder) sendSuccess() {
	RenderJSON(b.w, b.status, SuccessResponse{
		Data: b.data,
	}, b.headers)
}

func (b *ResponseBuilder) sendError() {
	RenderJSON(b.w, b.status, ErrorResponse{
		Error: errorInfo{
			Code:    b.errCode,
			Message: b.errMessage,
			Details: b.errDetails,
		},
	}, b.headers)
}
