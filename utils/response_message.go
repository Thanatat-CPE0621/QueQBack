package utils

import (
	"gitlab.com/paiduay/queq-hospital-api/config"
)

// ErrorMessagePrototype - a prototype for error message
type ErrorMessagePrototype struct {
	APIVersion string      `json:"apiVersion"`
	Error      errorObject `json:"error"`
}

// SuccessMessagePrototype -- a prototype for 	 success message
type SuccessMessagePrototype struct {
	APIVersion string     `json:"apiVersion"`
	Data       interface{} `json:"data"`
}

type errorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
func SuccessMessage(data interface{}) SuccessMessagePrototype {
	return SuccessMessagePrototype{APIVersion: config.Configs.APIVersion, Data: data}
}
