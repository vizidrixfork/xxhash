package xxhash

/*
#include "c-trunk/xxhash.h"
*/
import "C"

import (
	"hash"
	"unsafe"
)

type xxHash32 struct {
	seed  uint32
	sum   uint32
	state unsafe.Pointer
}

// Size returns the number of bytes Sum will return.
func (xx *xxHash32) Size() int {
	return 4
}

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
func (xx *xxHash32) BlockSize() int {
	return 8
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (xx *xxHash32) Sum(in []byte) []byte {
	s := xx.Sum32()
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

func (xx *xxHash32) Write(p []byte) (n int, err error) {
	switch {
	case xx.state == nil:
		return 0, ErrAlreadyComputed
	case len(p) > oneGb:
		return 0, ErrMemoryLimit
	}
	C.XXH32_update(xx.state, unsafe.Pointer(&p[0]), C.uint(len(p)))
	return len(p), nil
}

func (xx *xxHash32) Sum32() uint32 {
	if xx.state == nil {
		return xx.sum
	}
	xx.sum = uint32(C.XXH32_digest(xx.state))
	xx.state = nil
	return xx.sum
}

// Reset resets the Hash to its initial state.
func (xx *xxHash32) Reset() {
	if xx.state != nil {
		C.XXH32_digest(xx.state)
	}
	xx.state = C.XXH32_init(C.uint(xx.seed))
}

// NewS32 creates a new hash.Hash32 computing the 32bit xxHash checksum starting with the specific seed.
func NewS32(seed uint32) hash.Hash32 {
	h := &xxHash32{
		seed: seed,
	}
	h.Reset()
	return h
}

// New32 creates a new hash.Hash32 computing the 32bit xxHash checksum starting with the seed set to 0x0.
func New32() hash.Hash32 {
	return NewS32(0x0)
}

// Checksum32S returns the checksum of the input bytes with the specific seed.
func Checksum32S(in []byte, seed uint32) uint32 {
	return uint32(C.XXH32(unsafe.Pointer(&in[0]), C.uint(len(in)), C.uint(seed)))
}

// Checksum32 returns the checksum of the input data with the seed set to 0
func Checksum32(in []byte) uint32 {
	return Checksum32S(in, 0x0)
}
