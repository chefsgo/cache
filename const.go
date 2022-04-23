package cache

import "errors"

const (
	NAME = "cache"
)

var (
	errInvalidCacheConnection = errors.New("Invalid cache connection.")
)
