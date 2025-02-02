package gateway

import "errors"

// ErrNotFound is returned when a requested record is not found
var ErrNotFound = errors.New("ratings not found for a record")
