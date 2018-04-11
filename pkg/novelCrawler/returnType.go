package novelcrawler

// ErrorCode code of error type
type ErrorCode int

const (
	// Success success
	Success ErrorCode = iota
	// UnknownError error type
	UnknownError
	// RemoteNotExist remote url not exists
	RemoteNotExist
	// DBError the errors about database operation
	DBError
	// NovelName failed to get novel name
	NovelName
	// InternalError internal error
	InternalError
)

// ErrorString get the description of error code
func (code *ErrorCode) ErrorString() string {
	errors := [...]string{
		"Success",
		"Unknown error",
		"Remote novel not exists",
		"Database operation error",
		"Failed to get novel name",
		"Internal error",
	}
	if *code < Success || *code > InternalError {
		return "Wrong error code"
	}
	return errors[*code]
}
