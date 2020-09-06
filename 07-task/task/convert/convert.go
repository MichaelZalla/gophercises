// Package convert provides helper functions for translating between []byte and other common types
package convert

import (
	"encoding/binary"
)

// IntToBytes converts an integer type to a []byte value and returns it
func IntToBytes(v int) []byte {

	b := make([]byte, 8)

	binary.BigEndian.PutUint64(b, uint64(v))

	return b

}
