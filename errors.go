package gonfig

import "errors"

var (
	ErrExpectedNonNilOrEmpty = errors.New("provided value must not be nil or empty")
	ErrUnexpected            = errors.New("unexpected error")
	ErrMarshalling           = errors.New("unable to marshal config")
	ErrOverwriteRemove       = errors.New("failed to remove existing config file")
	ErrOverwriteDisabled     = errors.New("overwriting is disabled")
	ErrMustBeAddressable     = errors.New("provided variable must be a pointer")
)
