package gonfig

const (
	ErrExpectedNonNilOrEmpty = "provided value must not be nil or empty"
	ErrUnexpected            = "unexpected error"
	ErrOverwrite             = "unable to remove existing config file"
	ErrOverwriteDisabled     = "overwriting is disabled"
	ErrMarshalling           = "unable to marshal config"
	ErrCreatingConfig        = "unable to create config"
)
