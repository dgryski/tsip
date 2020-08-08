package tsip

import (
	"hash/fnv"
	"strconv"
	"testing"
)

var buf [8192]byte
var total uint64

func BenchmarkHash(b *testing.B)    { benchmarkHash(b, Hash) }
func BenchmarkHashASM(b *testing.B) { benchmarkHash(b, HashASM) }
func BenchmarkFNV(b *testing.B) {
	benchmarkHash(b, func(_ uint64, _ uint64, p []byte) uint64 {
		f := fnv.New64a()
		f.Write(p)
		return f.Sum64()
	})
}

func BenchmarkLevelDB(b *testing.B) {
	benchmarkHash(b, func(_ uint64, _ uint64, p []byte) uint64 {
		return uint64(leveldbHash(p))
	})
}

func leveldbHash(b []byte) uint32 {

	const (
		seed = 0xbc9f1d34
		m    = 0xc6a4a793
	)

	h := uint32(seed) ^ uint32(len(b)*m)

	for ; len(b) >= 4; b = b[4:] {
		h += uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
		h *= m
		h ^= h >> 16
	}
	switch len(b) {
	case 3:
		h += uint32(b[2]) << 16
		fallthrough
	case 2:
		h += uint32(b[1]) << 8
		fallthrough
	case 1:
		h += uint32(b[0])
		h *= m
		h ^= h >> 24
	}

	return h
}

func benchmarkHash(b *testing.B, h func(k0, k1 uint64, p []byte) uint64) {
	var sizes = []int{1, 2, 3, 4, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1024, 8192}
	for _, n := range sizes {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.SetBytes(int64(n))
			for i := 0; i < b.N; i++ {
				total += h(0, 0, buf[:n])
			}
		})
	}
}
