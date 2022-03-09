package tsip

import (
	"encoding/binary"
	"math/bits"
)

type sip struct {
	v0, v1 uint64
}

func (s *sip) round() {
	s.v0 += s.v1
	s.v1 = bits.RotateLeft64(s.v1, 13) ^ s.v0
	s.v0 = bits.RotateLeft64(s.v0, 35) + s.v1
	s.v1 = bits.RotateLeft64(s.v1, 17) ^ s.v0
	s.v0 = bits.RotateLeft64(s.v0, 21)
}

func Hash(k0, k1 uint64, p []byte) uint64 {

	s := sip{
		v0: k0 ^ 0x736f6d6570736575,
		v1: k1 ^ 0x646f72616e646f6d,
	}
	b := uint64(len(p)) << 56

	for len(p) >= 8 {
		m := binary.LittleEndian.Uint64(p[:8])
		s.v1 ^= m
		s.round()
		s.v0 ^= m
		p = p[8:]
	}

	switch len(p) {
	case 7:
		b |= uint64(p[6]) << 48
		fallthrough
	case 6:
		b |= uint64(p[5]) << 40
		fallthrough
	case 5:
		b |= uint64(p[4]) << 32
		fallthrough
	case 4:
		b |= uint64(p[3]) << 24
		fallthrough
	case 3:
		b |= uint64(p[2]) << 16
		fallthrough
	case 2:
		b |= uint64(p[1]) << 8
		fallthrough
	case 1:
		b |= uint64(p[0])
	}

	// last block
	s.v1 ^= b
	s.round()
	s.v0 ^= b

	// finalization
	s.v1 ^= 0xff
	s.round()
	s.v1 = bits.RotateLeft64(s.v1, 32)
	s.round()
	s.v1 = bits.RotateLeft64(s.v1, 32)

	return s.v0 ^ s.v1
}

type sip32 struct {
	v0, v1 uint32
}

func (s *sip32) round() {
	s.v0 += s.v1
	// rotation constants from chaskey core
	s.v1 = bits.RotateLeft32(s.v1, 5) ^ s.v0
	s.v0 = bits.RotateLeft32(s.v0, 8) + s.v1
	s.v1 = bits.RotateLeft32(s.v1, 13) ^ s.v0
	s.v0 = bits.RotateLeft32(s.v0, 7)
}

func Hash32(k0, k1 uint32, p []byte) uint32 {

	s := sip32{
		v0: k0 ^ 0x74656462,
		v1: k1 ^ 0x6c796765,
	}
	b := uint32(len(p)) << 24

	for len(p) >= 4 {
		m := binary.LittleEndian.Uint32(p[:4])
		s.v1 ^= m
		s.round()
		s.v0 ^= m
		p = p[4:]
	}

	switch len(p) {
	case 3:
		b |= uint32(p[2]) << 16
		fallthrough
	case 2:
		b |= uint32(p[1]) << 8
		fallthrough
	case 1:
		b |= uint32(p[0])
	}

	// last block
	s.v1 ^= b
	s.round()
	s.v0 ^= b

	// finalization
	s.v1 ^= 0xff
	s.round()
	s.v1 = bits.RotateLeft32(s.v1, 16)
	s.round()
	s.v1 = bits.RotateLeft32(s.v1, 16)
	//s.round()

	return s.v0 ^ s.v1
}
