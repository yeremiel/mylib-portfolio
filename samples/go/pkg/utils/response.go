package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yeremiel/mylib-portfolio/samples/go/pkg/i18n"
)

// RequestIDKey is the gin context key under which the request ID is stored.
const RequestIDKey = "requestId"

// APIResponse is the standard API response envelope.
// error: machine-readable error code (UPPER_SNAKE_CASE); used only in error responses.
// errors: per-field error details for 422 validation errors.
type APIResponse struct {
	Status   string       `json:"status"`
	Code     int          `json:"code"`
	Error    string       `json:"error,omitempty"`
	Message  string       `json:"message"`
	Data     interface{}  `json:"data,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
	Metadata Metadata     `json:"metadata"`
}

// FieldError is a field-level validation error (used in 422 responses).
type FieldError struct {
	Field   string `json:"field"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// Metadata is the response metadata.
type Metadata struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId,omitempty"`
}

// buildMetadata builds Metadata containing the timestamp and the request ID.
func buildMetadata(c *gin.Context) Metadata {
	m := Metadata{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if val, exists := c.Get(RequestIDKey); exists {
		if id, ok := val.(string); ok {
			m.RequestID = id
		}
	}
	return m
}

// SuccessResponse builds and sends a success response.
// messageKey: i18n message key (e.g. "auth.signup.success")
func SuccessResponse(c *gin.Context, code int, messageKey string, data interface{}) {
	locale := i18n.GetLocaleFromContext(c)
	translatedMessage := i18n.Translate(locale, messageKey)

	response := APIResponse{
		Status:   "success",
		Code:     code,
		Message:  translatedMessage,
		Data:     data,
		Metadata: buildMetadata(c),
	}
	c.JSON(code, response)
}

// ErrorResponseWithCode builds and sends an error response with a standard error code.
// errorCode: UPPER_SNAKE_CASE error code (e.g. "AUTH_INVALID_CREDENTIALS")
// messageKey: i18n message key or a plain message
func ErrorResponseWithCode(c *gin.Context, code int, errorCode string, messageKey string) {
	locale := i18n.GetLocaleFromContext(c)
	translatedMessage := i18n.Translate(locale, messageKey)

	response := APIResponse{
		Status:   "error",
		Code:     code,
		Error:    errorCode,
		Message:  translatedMessage,
		Metadata: buildMetadata(c),
	}
	c.JSON(code, response)
}

// ErrorResponseWithData builds and sends an error response with a standard error code
// and extra diagnostic data.
// data: additional information for the client (e.g. per-dependency status)
func ErrorResponseWithData(c *gin.Context, code int, errorCode string, messageKey string, data interface{}) {
	locale := i18n.GetLocaleFromContext(c)
	translatedMessage := i18n.Translate(locale, messageKey)

	response := APIResponse{
		Status:   "error",
		Code:     code,
		Error:    errorCode,
		Message:  translatedMessage,
		Data:     data,
		Metadata: buildMetadata(c),
	}
	c.JSON(code, response)
}

// ValidationErrorResponse builds and sends a 422 validation error response.
// messageKey: i18n message key
// fieldErrors: per-field error details
func ValidationErrorResponse(c *gin.Context, messageKey string, fieldErrors []FieldError) {
	locale := i18n.GetLocaleFromContext(c)
	translatedMessage := i18n.Translate(locale, messageKey)

	response := APIResponse{
		Status:   "error",
		Code:     http.StatusUnprocessableEntity,
		Error:    "REQUEST_VALIDATION_FAILED",
		Message:  translatedMessage,
		Errors:   fieldErrors,
		Metadata: buildMetadata(c),
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}
