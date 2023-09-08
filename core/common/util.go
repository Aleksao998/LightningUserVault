package common

import "encoding/binary"

// Int64ToBytes converts int64 to []byte using little endian format.
func Int64ToBytes(i int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))

	return bytes
}

// BytesToInt64 converts []byte to int64 using little endian format.
func BytesToInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}
