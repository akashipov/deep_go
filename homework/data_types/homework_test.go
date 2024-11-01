package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func ToLittleEndianV4WithoutIf(number uint32) (res uint32) {
	const (
		byteLength = 1 << 3
		mask       = uint32(byteLength<<5 - 1)
	)

	p := mask & number

	for i := 24; i >= 0; i -= byteLength {
		res |= p << i
		number >>= byteLength
		p = mask & number
	}

	return
}

func ToLittleEndian[T uint64 | uint16 | uint32](number T) (reversedBytes T) {
	size := int(unsafe.Sizeof(number))
	pointer := unsafe.Pointer(&number)

	k := number

	for _ = range size {
		reversedBytes = (reversedBytes << 8) | T(*(*byte)(pointer))
		pointer = unsafe.Add(pointer, 1)
		k <<= 8
		if k == 0 {
			return
		}
	}
	return
}

func ToLittleEndianUnsafe(number uint32) (res uint32) {
	p := unsafe.Pointer(&number)

	k := number

	for i := 24; i >= 0; i -= 8 {
		fmt.Printf("%08x\n", *(*uint8)(p))
		res |= uint32(*(*uint8)(p)) << i
		k >>= 8
		if k == 0 {
			return
		}
		p = unsafe.Pointer(uintptr(p) + 1)
	}

	return
}

func ToLittleEndianV1WithIf(number uint32) (res uint32) {
	const (
		byteLength = 1 << 3
	)

	p := byte(number)

	for i := 24; i >= 0; i -= byteLength {
		res |= uint32(p) << i
		number >>= byteLength
		if number == 0 {
			return
		}
		p = byte(number)
	}

	return
}

func ToLittleEndianV1MaskWithIf(number uint32) (res uint32) {
	const (
		byteLength = 1 << 3
		mask       = uint32(byteLength<<5 - 1)
	)

	p := mask & number

	for i := 24; i >= 0; i -= byteLength {
		res |= p << i
		number >>= byteLength
		if number == 0 {
			return
		}
		p = mask & number
	}

	return
}

func BenchmarkBigEndian(b *testing.B) {
	length := 1000
	cases := make([]uint32, 0, length)
	for i := range length {
		cases = append(cases, uint32(i))
	}

	b.ResetTimer()

	for _ = range b.N {
		for _, c := range cases {
			ToLittleEndian(c)
		}
	}
}

func BenchmarkBigEndianMaskWithIf(b *testing.B) {
	length := 1000
	cases := make([]uint32, 0, length)
	for i := range length {
		cases = append(cases, uint32(i))
	}

	b.ResetTimer()

	for _ = range b.N {
		for _, c := range cases {
			ToLittleEndianV1MaskWithIf(c)
		}
	}
}

func ToLittleEndianLast[T uint32](number T) T {
	size := int(unsafe.Sizeof(number))

	ptr := (*[8]byte)(unsafe.Pointer(&number))

	var result T
	for i := 0; i < size; i++ {
		result |= T(ptr[i]) << (8 * (size - 1 - i))
	}

	return result
}

func BenchmarkBigEndianUnsafe(b *testing.B) {
	length := 1000
	cases := make([]uint32, 0, length)
	for i := range length {
		cases = append(cases, uint32(i))
	}

	b.ResetTimer()

	for _ = range b.N {
		for _, c := range cases {
			ToLittleEndianUnsafe(c)
		}
	}
}

func BenchmarkBigEndianV1WithIf(b *testing.B) {
	length := 1000
	cases := make([]uint32, 0, length)
	for i := range length {
		cases = append(cases, uint32(i))
	}

	b.ResetTimer()

	for _ = range b.N {
		for _, c := range cases {
			ToLittleEndianLast(c)
		}
	}
}

func TestSerializationProperties(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"test case #6": {
			number: 0x12345678,
			result: 0x78563412,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianLast(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
