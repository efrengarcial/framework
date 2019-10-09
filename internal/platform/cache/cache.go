// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/oogway/go-cache

package cache

import (
	"time"
)

// Length of time to cache an item.
const (
	DefaultExpiryTime  = time.Duration(0)
	ForEverNeverExpiry = time.Duration(-1)
)

// Getter is an interface for getting / decoding an element from a cache.
type Getter interface {
	// Get the content associated with the given key. decoding it into the given
	// pointer.
	//
	// Returns:
	//   - nil if the value was successfully retrieved and ptrValue set
	//   - ErrCacheMiss if the value was not in the cache
	//   - an implementation specific error otherwise
	Get(key string, ptrValue interface{}) error
}

// Cache is an interface to an expiring cache.  It behaves (and is modeled) like
// the Memcached interface.  It is keyed by strings (250 bytes at most).
//
// Many callers will make exclusive use of Set and Get, but more exotic
// functions are also available.
//
// Example
//
// Here is a typical Get/Set interaction:
//
//   var items []*Item
//   if err := cache.Get("items", &items); err != nil {
//     items = loadItems()
//     go cache.Set("items", items, cache.DefaultExpiryTime)
//   }
//
// Note that the caller will frequently not wait for Set() to complete.
//
// Errors
//
// It is assumed that callers will infrequently check returned errors, since any
// request should be fulfillable without finding anything in the cache.  As a
// result, all errors other than ErrCacheMiss and ErrNotStored will be logged to
// revel.ERROR, so that the developer does not need to check the return value to
// discover things like deserialization or connection errors.
type Cache interface {
	// The Cache implements a Getter.
	Getter

	// Set the given key/value in the cache, overwriting any existing value
	// associated with that key.  Keys may be at most 250 bytes in length.
	//
	// Returns:
	//   - nil on success
	//   - an implementation specific error otherwise
	Set(key string, value interface{}, expires time.Duration) error

	// SetFields will atomically set a field of a Hash.
	SetFields(key string, value map[string]interface{}, expires time.Duration) error

	// Get the content associated multiple keys at once.  On success, the caller
	// may decode the values one at a time from the returned Getter.
	//
	// Returns:
	//   - the value getter, and a nil error if the operation completed.
	//   - an implementation specific error otherwise
	GetMulti(keys ...string) (Getter, error)

	// Delete the given key from the cache.
	//
	// Returns:
	//   - nil on a successful delete
	//   - ErrCacheMiss if the value was not in the cache
	//   - an implementation specific error otherwise
	Delete(key string) error

	// Add the given key/value to the cache ONLY IF the key does not already exist.
	//
	// Returns:
	//   - nil if the value was added to the cache
	//   - ErrNotStored if the key was already present in the cache
	//   - an implementation-specific error otherwise
	Add(key string, value interface{}, expires time.Duration) error

	// Set the given key/value in the cache ONLY IF the key already exists.
	//
	// Returns:
	//   - nil if the value was replaced
	//   - ErrNotStored if the key does not exist in the cache
	//   - an implementation specific error otherwise
	Replace(key string, value interface{}, expires time.Duration) error

	// Expire all cache entries immediately.
	// This is not implemented for the memcached cache (intentionally).
	// Returns an implementation specific error if the operation failed.
	Flush() error

	// Get all currently set keys. This can be super slow so use with care.
	Keys() ([]string, error)
}