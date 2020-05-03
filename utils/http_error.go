package utils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"strings"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

// Codes for error Identification (TODO: Add and use further codes)
const (
	ConstFailedValidation = "40001" // error code for a failed validation
)

// Keys for ErrorType
const (
	ErrorTypeBadRequest       = "BadRequest"
	ErrorTypeBinding          = "Binding"
	ErrorTypeValidation       = "Validation"
	ErrorTypeResourceNotFound = "ResourceNotFound"
	ErrorTypeURLNotFound      = "URLNotFound"
	ErrorTypeDatabase         = "Database"
	ErrorTypeInternalServer   = "InternalServer"
	ErrorTypeBadGateway       = "BadGateway"
	ErrorTypeUnauthorized     = "Unauthorized"
	ErrorTypeForbidden        = "Forbidden"
	ErrorTypeConflict         = "Conflict"
	ErrorTypeTooManyRequests  = "TooManyRequests"
)

// HasHTTPStatus Error Interface which contains an HTTP Status and a specific error type
type HasHTTPStatus interface {
	GetHTTPStatusCode() int
	GetErrorType() string
}

// CommonError Parent struct for errors which are used in microservices
// Each child error contains an error type, http status code and the previous 'thrown' error
// It implements the HasHTTPStatus Interface
type CommonError struct {
	Err   error
	Code  int
	XType string
}

func (e *CommonError) Error() string {
	return e.Err.Error()
}
func (e *CommonError) GetHTTPStatusCode() int {
	return e.Code
}
func (e *CommonError) GetErrorType() string {
	return e.XType
}

// ErrorBadRequest Error for 400 Responses
type ErrorBadRequest struct {
	CommonError
}

// ErrorUnauthorized Error for 401 Responses
type ErrorUnauthorized struct {
	CommonError
}

// ErrorBinding Error for 400 Responses when echo.Bind fails
type ErrorBinding struct {
	CommonError
}

// ErrorValidation Error for 401 Responses when validation fails
type ErrorValidation struct {
	CommonError
}

// ErrorDatabase Error for 500 Responses when database functions fails
type ErrorDatabase struct {
	CommonError
}

// ErrorResourceNotFound Error for 404 Responses when a resource is not found
type ErrorResourceNotFound struct {
	CommonError
}

// ErrorURLNotFound Error for 404 Responses when a URL is not available
type ErrorURLNotFound struct {
	CommonError
}

// ErrorInternalServer Error for 500 Responses when server runs into an error.
// Also used for errors, which are not clearly specified
type ErrorInternalServer struct {
	CommonError
}

// ErrorBadGateway Error for 502 Responses when request can not be served.
type ErrorBadGateway struct {
	CommonError
}

// ErrorForbidden Error for 403 Responses when access to the requested resource is forbidden.
type ErrorForbidden struct {
	CommonError
}

// ErrorConflict Error for 409 Responses when the request conflicts with the current state of the server.
type ErrorConflict struct {
	CommonError
}

// ErrorTooManyRequests Error for 429 Responses
type ErrorTooManyRequests struct {
	CommonError
}

func Errorf(xtype string, message string, args ...interface{}) error {
	return Error(fmt.Sprintf(message, args...), xtype)
}

// Error Factory method for creating CommonErrors
func Error(i interface{}, xtype string) error {
	var err error
	if str, ok := i.(string); ok {
		err = errors.New(str)
	} else if commonError, ok := i.(HasHTTPStatus); ok {
		return commonError.(error)
	} else if err, ok = i.(error); !ok {
		panic(fmt.Sprintf("I don't know how to handle that type: %v", i))
	}

	switch xtype {
	case ErrorTypeBadRequest:
		return &ErrorBadRequest{CommonError{err, http.StatusBadRequest, xtype}}
	case ErrorTypeBinding:
		return &ErrorBinding{CommonError{err, http.StatusBadRequest, xtype}}
	case ErrorTypeValidation:
		return &ErrorValidation{CommonError{err, http.StatusBadRequest, xtype}}
	case ErrorTypeDatabase:
		return &ErrorDatabase{CommonError{err, http.StatusInternalServerError, xtype}}
	case ErrorTypeResourceNotFound:
		return &ErrorResourceNotFound{CommonError{err, http.StatusNotFound, xtype}}
	case ErrorTypeURLNotFound:
		return &ErrorURLNotFound{CommonError{err, http.StatusNotFound, xtype}}
	case ErrorTypeUnauthorized:
		return &ErrorUnauthorized{CommonError{err, http.StatusUnauthorized, xtype}}
	case ErrorTypeBadGateway:
		return &ErrorBadGateway{CommonError{err, http.StatusBadGateway, xtype}}
	case ErrorTypeForbidden:
		return &ErrorForbidden{CommonError{err, http.StatusForbidden, xtype}}
	case ErrorTypeConflict:
		return &ErrorConflict{CommonError{err, http.StatusConflict, xtype}}
	case ErrorTypeTooManyRequests:
		return &ErrorTooManyRequests{CommonError{err, http.StatusTooManyRequests, xtype}}
	default:
		return &ErrorInternalServer{CommonError{err, http.StatusInternalServerError, xtype}}
	}
}

// ValidationErrorStructure - Representation of a validation error
type ValidationErrorStructure struct {
	Class     string `json:"class,omitempty"`
	Field     string `json:"field"`
	Validator string `json:"validator"`
	Message   string `json:"message,omitempty"`
}

// HTTPError - Basic Implementation of ClientError
type HTTPError struct {
	Code             string                     `json:"code,omitempty"`
	Type             string                     `json:"type"`
	Message          string                     `json:"message"`
	MessageID        string                     `json:"messageId,omitempty"`
	ValidationErrors []ValidationErrorStructure `json:"validationErrors,omitempty"`
	Status           int                        `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// HTTPErrorHandler Error handler to use in echo's middleware
// Example:
// 			router := echo.New()
// 			router.HTTPErrorHandler = commons.HTTPErrorHandler
func HTTPErrorHandler(err error, c echo.Context) {

	//TODO: log resultErr to external monitoring service

	// Because echo.NotFoundHandler return an echo.HTTPError and to avoid duplicate of code at router initialization
	// the error given by echo will be overwritten by an error of our own type
	if (reflect.TypeOf(err) == reflect.TypeOf(&echo.HTTPError{})) && (err.(*echo.HTTPError).Code == http.StatusNotFound) {
		err = Error(err, ErrorTypeURLNotFound)
	}

	// Because echo.MethodNotAllowedHandler return an echo.HTTPError and to avoid duplicate of code at router initialization
	// the error given by echo will be overwritten by an error of our own type
	if (reflect.TypeOf(err) == reflect.TypeOf(&echo.HTTPError{})) && (err.(*echo.HTTPError).Code == http.StatusUnauthorized) {
		err = Error(err, ErrorTypeUnauthorized)
	}

	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	switch err.(type) {
	case *ErrorValidation:
		err = c.JSON(err.(*ErrorValidation).Code, validationErrorToHTTPError(err.(*ErrorValidation), requestID))
	case HasHTTPStatus:
		httpError := newHTTPError(err.(HasHTTPStatus), nil, requestID)
		err = c.JSON(err.(HasHTTPStatus).GetHTTPStatusCode(), httpError)
	default:
		httpError := newHTTPError(Error(err, ErrorTypeInternalServer).(HasHTTPStatus), nil, requestID)
		err = c.JSON(http.StatusInternalServerError, httpError)
	}
}

// newHTTPError - Create a new HTTPError instance
func newHTTPError(err HasHTTPStatus, valErrors []ValidationErrorStructure, requestID string) error {
	errMsg := err.(error).Error()

	return &HTTPError{
		Code:             "",
		Type:             err.GetErrorType(),
		Message:          errMsg,
		MessageID:        requestID,
		ValidationErrors: valErrors,
		Status:           err.GetHTTPStatusCode(),
	}
}

// validationErrorToHTTPError - Generate a HTTPError for a bad request including validation errors
func validationErrorToHTTPError(err *ErrorValidation, requestID string) error {

	// extract validation error information from validation structs
	valErrors := []ValidationErrorStructure{}
	if verr, ok := err.Err.(validator.ValidationErrors); ok {
		for i := 0; i < len(verr); i++ {
			structNames := strings.Split(verr[i].StructNamespace(), ".")
			structNS := structNames[0]
			valError := ValidationErrorStructure{structNS, verr[i].StructField(), verr[i].ActualTag(), verr[i].Translate(nil)}
			valErrors = append(valErrors, valError)
		}
	}

	// and paste it to http error struct
	httpError := newHTTPError(err, valErrors, requestID)
	return httpError
}

// DetermineErrorTypeByStatusCode Determines a suitable ErrorType based on the status code
func DetermineErrorTypeByStatusCode(statusCode int) string {
	switch statusCode {
	case 400:
		return ErrorTypeBadRequest
	case 401:
		return ErrorTypeUnauthorized
	case 403:
		return ErrorTypeForbidden
	case 404:
		return ErrorTypeResourceNotFound
	case 409:
		return ErrorTypeConflict
	case 429:
		return ErrorTypeTooManyRequests
	case 500:
		return ErrorTypeInternalServer
	default:
		return ErrorTypeInternalServer
	}
}
