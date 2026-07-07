package repository

import "errors"

// ErrNotFound is returned by repositories when an entity does not exist.
// Use errors.Is to detect it across error wrapping.
var ErrNotFound = errors.New("entity not found")

// ErrConflict is returned by repositories when a unique constraint
// (e.g. email, username) would be violated.
var ErrConflict = errors.New("entity conflict")
