package bloomFilter

import (
	"log"
	"math"
	"unsafe"
)

// helper
var mask = []uint8{1, 2, 4, 8, 16, 32, 64, 128}

func getSize(ui64 uint64) (size uint64, exponent uint64) {
	if ui64 < uint64(512) {
		ui64 = uint64(512)
	}
	size = uint64(1)
	for size < ui64 {
		size <<= 1
		exponent++
	}
	return size, exponent
}

func calcSizeByWrongPositives(numEntries, wrongs float64) (uint64, uint64) {
	size := -1 * numEntries * math.Log(wrongs) / math.Pow(float64(0.69314718056), 2)
	locs := math.Ceil(float64(0.69314718056) * size / numEntries)
	return uint64(size), uint64(locs)
}

func New(params ...float64) (bloomfilter Bloom) {
	var entries, locs uint64
	if len(params) == 2 {
		if params[1] < 1 {
			entries, locs = calcSizeByWrongPositives(params[0], params[1])
		} else {
			entries, locs = uint64(params[0]), uint64(params[1])
		}
	} else {
		log.Fatal("Bad usage! Please check Readme.md")
	}
	size, exponent := getSize(uint64(entries))
	bloomfilter = Bloom{
		sizeExp: exponent,
		size:    size - 1,
		setLocs: locs,
		shift:   64 - exponent,
	}
	bloomfilter.Size(size)
	return bloomfilter
}


type Bloom struct {
	bitset  Bitset
	sizeExp uint64
	size    uint64
	setLocs uint64
	shift   uint64
}

func (bl Bloom) Add(entry []byte) {
	l, h := bl.sipHash(entry)
	for i := uint64(0); i < bl.setLocs; i++ {
		bl.bitset.Set((h + i*l) & bl.size)
	}
}

func (bl Bloom) Has(entry []byte) bool {
	l, h := bl.sipHash(entry)
	for i := uint64(0); i < bl.setLocs; i++ {
		switch bl.bitset.IsSet((h + i*l) & bl.size) {
		case false:
			return false
		}
	}
	return true
}

func (bl *Bloom) Size(sz uint64) {
	(*bl).bitset.Size(sz)
}

// Clear
// resets the Bloom filter
func (bl *Bloom) Clear() {
	(*bl).bitset.Clear()
}


type Bitset struct {
	bs    []uint64
	start uintptr
}

func (bs *Bitset) Size(sz uint64) {
	(*bs).bs = make([]uint64, sz>>6)
	(*bs).start = uintptr(unsafe.Pointer(&bs.bs[0]))
}

func (bs *Bitset) Set(idx uint64) {
	ptr := unsafe.Pointer(bs.start + uintptr(idx>>3))
	*(*uint8)(ptr) |= mask[idx%8]
}

func (bs *Bitset) Clear() {
	for i, _ := range (*bs).bs {
		(*bs).bs[i] = 0
	}
}

func (bs *Bitset) IsSet(idx uint64) bool {
	ptr := unsafe.Pointer(bs.start + uintptr(idx>>3))
	r := ((*(*uint8)(ptr)) >> (idx % 8)) & 1
	return r == 1
}
