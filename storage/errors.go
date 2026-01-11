package storage

import "errors"

// Drive-related errors
var (
	ErrDriveNotFound            = errors.New("drive not found")
	ErrDriveNotFoundOrInUse     = errors.New("drive not found or already in use")
	ErrAlreadyAdopted           = errors.New("drive already adopted")
	ErrNoDrivesToRemove         = errors.New("no drives to remove")
	ErrDuplicateDriveKey        = errors.New("duplicate drive key found")
)

// Pool-related errors
var (
	ErrPoolNotFound        = errors.New("pool not found")
	ErrPoolInUse           = errors.New("pool is currently in use")
	ErrPoolAlreadyExists   = errors.New("pool with the same UUID already exists")
	ErrPoolNotOffline      = errors.New("cannot delete a pool that is not offline")
	ErrPoolFormatRequired  = errors.New("pool format must be specified")
	ErrInsufficientDrives  = errors.New("insufficient drives for the requested pool type")
	ErrInvalidPoolType     = errors.New("invalid pool type")
	ErrUnsupportedFormat   = errors.New("unsupported pool format")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrUuidTooShort        = errors.New("uuid length is less than requested length")
	ErrPoolNotInMemory     = errors.New("pool not found in memory")
)

// Generic errors
var (
	ErrNotFound = errors.New("resource not found")
)
