package storage

import "errors"

// Drive-related errors
var (
	ErrAlreadyAdopted       = errors.New("drive already adopted")
	ErrDriveNotFound        = errors.New("drive not found")
	ErrDriveNotFoundOrInUse = errors.New("drive not found or already in use")
	ErrDuplicateDriveKey    = errors.New("duplicate drive key found")
	ErrNoDrivesToRemove     = errors.New("no drives to remove")
)

// Pool-related errors
var (
	ErrInsufficientDrives = errors.New("insufficient drives for the requested pool type")
	ErrInvalidPoolType    = errors.New("invalid pool type")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrPoolAlreadyExists  = errors.New("pool with the same UUID already exists")
	ErrPoolFormatRequired = errors.New("pool format must be specified")
	ErrPoolInUse          = errors.New("pool is currently in use")
	ErrPoolNotFound       = errors.New("pool not found")
	ErrPoolNotInMemory    = errors.New("pool not found in memory")
	ErrPoolNotOffline     = errors.New("cannot delete a pool that is not offline")
	ErrUnsupportedFormat  = errors.New("unsupported pool format")
	ErrUuidTooShort       = errors.New("uuid length is less than requested length")
)

// Generic errors
var (
	ErrNotFound = errors.New("resource not found")
)
