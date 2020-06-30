package bloompass

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/spaolacci/murmur3"
	"hash"
	"hash/fnv"
)

const EXIST_NO = 0
const EXIST_YES = 1

type Bloom struct {
	maps    map[string]*roaring.Bitmap
	ciphers map[string]hash.Hash32
}

func NewBloom(seed uint32) *Bloom {
	maps := make(map[string]*roaring.Bitmap, 0)
	ciphers := make(map[string]hash.Hash32, 0)
	ciphers["murmur3"] = murmur3.New32WithSeed(seed)
	ciphers["fnv"] = fnv.New32a()
	for k, _ := range ciphers {
		maps[k] = roaring.New()
	}
	return &Bloom{
		maps:    maps,
		ciphers: ciphers,
	}
}

func (b *Bloom) Add(s string) {
	for name, h := range b.ciphers {
		h.Reset()
		h.Write([]byte(s))
		idx := h.Sum32()
		b.maps[name].Add(idx)
	}
}

func (b *Bloom) Exists(s string) int {
	result := 0
	for name, bitmap := range b.maps {
		h := b.ciphers[name]
		h.Reset()
		h.Write([]byte(s))
		idx := h.Sum32()

		if bitmap.Contains(idx) {
			result += 1
		}
	}
	if result == 0 {
		return EXIST_NO
	}
	if result == len(b.ciphers) {
		return EXIST_YES
	}

	// we consider partial matching a NO, a possible collision
	return EXIST_NO
}
