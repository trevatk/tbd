package authoritative

import "errors"

var (
	errKeyNotFound = errors.New("key not found")
	errKeyExists   = errors.New("key exists")

	errInvalidRequestID = errors.New("invalid request id")
)
