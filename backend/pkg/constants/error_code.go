package constants

type (
	ErrorCode = string
)

const (
	ErrCodeSuccess     ErrorCode = "success"
	ErrCodeUnknown     ErrorCode = "unknown"
	ErrCodeInvalidArgs ErrorCode = "invalid_arguments"
	ErrCodeNotFound    ErrorCode = "not_found"

	ErrFileCheckSumInvalid ErrorCode = "file_checksum_invalid"
)
