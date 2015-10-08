// https://en.wikipedia.org/wiki/SipHash
// created by Jean-Philippe Aumasson and Daniel J. Bernstein.

package bloomFilter

func (bl Bloom) sipHash(p []byte) (l, h uint64) {
	// Initialization.
	v0 := uint64(8317987320269560794) // k0 ^ 0x736f6d6570736575
	v1 := uint64(7237128889637516672) // k1 ^ 0x646f72616e646f6d
	v2 := uint64(7816392314733513934) // k0 ^ 0x6c7967656e657261
	v3 := uint64(8387220255325274014) // k1 ^ 0x7465646279746573
	t := uint64(len(p)) << 56
	for len(p) >= 8 {
		m := uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 |
		uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56
		v3 ^= m
		//R1.
		v0 += v1
		v1 = v1<<13 | v1>>51
		v1 ^= v0
		v0 = v0<<32 | v0>>32
		v2 += v3
		v3 = v3<<16 | v3>>48
		v3 ^= v2
		v0 += v3
		v3 = v3<<21 | v3>>43
		v3 ^= v0
		v2 += v1
		v1 = v1<<17 | v1>>47
		v1 ^= v2
		v2 = v2<<32 | v2>>32
		//R2
		v0 += v1
		v1 = v1<<13 | v1>>51
		v1 ^= v0
		v0 = v0<<32 | v0>>32
		v2 += v3
		v3 = v3<<16 | v3>>48
		v3 ^= v2
		v0 += v3
		v3 = v3<<21 | v3>>43
		v3 ^= v0
		//R3
		v2 += v1
		v1 = v1<<17 | v1>>47
		v1 ^= v2
		v2 = v2<<32 | v2>>32
		v0 ^= m
		p = p[8:]
	}
	switch len(p) {
	case 7:
		t |= uint64(p[6]) << 48
		fallthrough
	case 6:
		t |= uint64(p[5]) << 40
		fallthrough
	case 5:
		t |= uint64(p[4]) << 32
		fallthrough
	case 4:
		t |= uint64(p[3]) << 24
		fallthrough
	case 3:
		t |= uint64(p[2]) << 16
		fallthrough
	case 2:
		t |= uint64(p[1]) << 8
		fallthrough
	case 1:
		t |= uint64(p[0])
	}
	v3 ^= t
	//R1
	v0 += v1
	v1 = v1<<13 | v1>>51
	v1 ^= v0
	v0 = v0<<32 | v0>>32
	v2 += v3
	v3 = v3<<16 | v3>>48
	v3 ^= v2
	v0 += v3
	v3 = v3<<21 | v3>>43
	v3 ^= v0
	v2 += v1
	v1 = v1<<17 | v1>>47
	v1 ^= v2
	v2 = v2<<32 | v2>>32
	//R2
	v0 += v1
	v1 = v1<<13 | v1>>51
	v1 ^= v0
	v0 = v0<<32 | v0>>32
	v2 += v3
	v3 = v3<<16 | v3>>48
	v3 ^= v2
	v0 += v3
	v3 = v3<<21 | v3>>43
	v3 ^= v0
	v2 += v1
	v1 = v1<<17 | v1>>47
	v1 ^= v2
	v2 = v2<<32 | v2>>32
	v0 ^= t
	//F
	v2 ^= 0xff
	//R1
	v0 += v1
	v1 = v1<<13 | v1>>51
	v1 ^= v0
	v0 = v0<<32 | v0>>32
	v2 += v3
	v3 = v3<<16 | v3>>48
	v3 ^= v2
	v0 += v3
	v3 = v3<<21 | v3>>43
	v3 ^= v0
	v2 += v1
	v1 = v1<<17 | v1>>47
	v1 ^= v2
	v2 = v2<<32 | v2>>32
	//R2
	v0 += v1
	v1 = v1<<13 | v1>>51
	v1 ^= v0
	v0 = v0<<32 | v0>>32
	v2 += v3
	v3 = v3<<16 | v3>>48
	v3 ^= v2
	v0 += v3
	v3 = v3<<21 | v3>>43
	v3 ^= v0
	v2 += v1
	v1 = v1<<17 | v1>>47
	v1 ^= v2
	v2 = v2<<32 | v2>>32
	//R3
	v0 += v1
	v1 = v1<<13 | v1>>51
	v1 ^= v0
	v0 = v0<<32 | v0>>32
	v2 += v3
	v3 = v3<<16 | v3>>48
	v3 ^= v2
	v0 += v3
	v3 = v3<<21 | v3>>43
	v3 ^= v0
	v2 += v1
	v1 = v1<<17 | v1>>47
	v1 ^= v2
	v2 = v2<<32 | v2>>32
	//R4
	v0 += v1
	v1 = v1<<13 | v1>>51
	v1 ^= v0
	v0 = v0<<32 | v0>>32
	v2 += v3
	v3 = v3<<16 | v3>>48
	v3 ^= v2
	v0 += v3
	v3 = v3<<21 | v3>>43
	v3 ^= v0
	v2 += v1
	v1 = v1<<17 | v1>>47
	v1 ^= v2
	v2 = v2<<32 | v2>>32

	hash := v0 ^ v1 ^ v2 ^ v3
	h = hash >> bl.shift
	l = hash << bl.shift >> bl.shift
	return l, h
}
