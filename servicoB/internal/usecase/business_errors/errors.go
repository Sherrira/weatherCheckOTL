package business_errors

import "errors"

var ErrCepNotFound = errors.New("can not find zipcode")
var ErrCepValidationFailed = errors.New("invalid zipcode")
var ErrFetchTemperatureFailed = errors.New("can not fetch temperature")
