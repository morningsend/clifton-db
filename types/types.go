// Common types used across different layers in the database.
package types

// Key value store
type Key []byte
type Value []byte

// hashing
type HashCode uint32
type HashFunc func(Key) HashCode


