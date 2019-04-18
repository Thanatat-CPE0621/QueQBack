package utils

import (
	"gitlab.com/paiduay/queq-hospital-api/config"
)

const (
	AuthErrorMessage = ""
	// AuthRequiredError is a message sending when there is no Authorization field in the header
	AuthRequiredError = "authorization required"
	TokenError        = "invalid token"
	PermissionError   = "permission denied"
)

// ErrorMessagePrototype - a prototype for error message
type ErrorMessagePrototype struct {
	APIVersion string      `json:"apiVersion"`
	Error      errorObject `json:"error"`
}

// SuccessMessagePrototype -- a prototype for success message
type SuccessMessagePrototype struct {
	APIVersion string     `json:"apiVersion"`
	Data       DataObject `json:"data"`
}

type errorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DataObject struct {
	Kind        *string     `json:"kind,omitempty"`
	Title       *string     `json:"title,omitempty"`
	Description *string     `json:"description,omitempty"`
	Item        interface{} `json:"item,omitempty"`
	Items       interface{} `json:"items,omitempty"`
}

// ErrorMessage - return an error message
func ErrorMessage(message string, code int) ErrorMessagePrototype {
	err := errorObject{
		Code:    code,
		Message: message,
	}

	return ErrorMessagePrototype{APIVersion: config.Configs.APIVersion, Error: err}
}

// SuccessMessage - return an success message
func SuccessMessage(data DataObject) SuccessMessagePrototype {
	return SuccessMessagePrototype{APIVersion: config.Configs.APIVersion, Data: data}
}
