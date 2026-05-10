package app

import "errors"

var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrRouterNotInitialized   = errors.New("router not initialized")
)