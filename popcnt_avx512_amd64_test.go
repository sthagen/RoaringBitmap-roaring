//go:build amd64 && !appengine
// +build amd64,!appengine

package roaring

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// edge lengths exercise the AVX-512 main loop (multiples of 32) and the scalar
// POPCNTQ tail (len % 32 != 0), including the empty and sub-block cases.
var avx512TestLengths = []int{0, 1, 2, 3, 4, 5, 7, 8, 15, 16, 17, 31, 32, 33, 63, 64, 65, 1023, 1024, 1025}

func TestPopcntAVX512MatchesGo(t *testing.T) {
	if !useAVX512Popcnt {
		t.Skip("AVX512_VPOPCNTDQ not available")
	}
	r := rand.New(rand.NewSource(42))
	for _, n := range avx512TestLengths {
		s := randomUint64Slice(r, n)
		m := randomUint64Slice(r, n)

		assert.Equalf(t, popcntSliceGo(s), popcntSliceAVX512(s), "popcntSlice len=%d", n)
		assert.Equalf(t, popcntMaskSliceGo(s, m), popcntMaskSliceAVX512(s, m), "popcntMaskSlice len=%d", n)
		assert.Equalf(t, popcntAndSliceGo(s, m), popcntAndSliceAVX512(s, m), "popcntAndSlice len=%d", n)
		assert.Equalf(t, popcntOrSliceGo(s, m), popcntOrSliceAVX512(s, m), "popcntOrSlice len=%d", n)
		assert.Equalf(t, popcntXorSliceGo(s, m), popcntXorSliceAVX512(s, m), "popcntXorSlice len=%d", n)
	}
}

// The dispatchers must agree with the portable implementations whatever the
// running CPU supports.
func TestPopcntDispatchMatchesGo(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for _, n := range avx512TestLengths {
		s := randomUint64Slice(r, n)
		m := randomUint64Slice(r, n)

		assert.Equalf(t, popcntSliceGo(s), popcntSlice(s), "popcntSlice len=%d", n)
		assert.Equalf(t, popcntMaskSliceGo(s, m), popcntMaskSlice(s, m), "popcntMaskSlice len=%d", n)
		assert.Equalf(t, popcntAndSliceGo(s, m), popcntAndSlice(s, m), "popcntAndSlice len=%d", n)
		assert.Equalf(t, popcntOrSliceGo(s, m), popcntOrSlice(s, m), "popcntOrSlice len=%d", n)
		assert.Equalf(t, popcntXorSliceGo(s, m), popcntXorSlice(s, m), "popcntXorSlice len=%d", n)
	}
}

func BenchmarkPopcntSliceAVX512(b *testing.B) {
	if !useAVX512Popcnt {
		b.Skip("AVX512_VPOPCNTDQ not available")
	}
	r := rand.New(rand.NewSource(1))
	s := randomUint64Slice(r, 1024)
	b.SetBytes(int64(len(s) * 8))
	var sink uint64
	for b.Loop() {
		sink = popcntSliceAVX512(s)
	}
	_ = sink
}

func BenchmarkPopcntAndSliceAVX512(b *testing.B) {
	if !useAVX512Popcnt {
		b.Skip("AVX512_VPOPCNTDQ not available")
	}
	r := rand.New(rand.NewSource(1))
	s := randomUint64Slice(r, 1024)
	m := randomUint64Slice(r, 1024)
	b.SetBytes(int64(len(s) * 8))
	var sink uint64
	for b.Loop() {
		sink = popcntAndSliceAVX512(s, m)
	}
	_ = sink
}
