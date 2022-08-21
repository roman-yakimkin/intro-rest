package errors

import "errors"

var ErrEntityNotFound = errors.New("entity not found")
var ErrOnEntityDeleting = errors.New("error on entity deleting")
