package bloomFilter

import (
	"math"
	"unsafe"
	"os"
	"bufio"
	"bytes"
	"io"
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

func New(filterSize, failRate float64) (filter *Bloom, err error) {
	entries, locs := calcSizeByWrongPositives(filterSize, failRate)
	size, exponent := getSize(uint64(entries))
	filter = &Bloom{
		sizeExp: exponent,
		size:    size - 1,
		setLocs: locs,
		shift:   64 - exponent,
	}
	filter.Size(size)
	return
}

// Creates a Bloom filter from a os file
func NewFromFile(filePath string, failRate float64) (filter *Bloom, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	return NewFromReadSeeker(file, failRate)
}

func NewFromReadSeeker(reader io.ReadSeeker, failRate float64) (filter *Bloom, err error) {
	filterSize, err := lineCounter(reader)
	if err != nil {
		return
	}
	filter, err = New(float64(filterSize), failRate)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		data := scanner.Bytes()
		filter.Add(bytes.TrimSpace(bytes.ToLower(data)))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return
}

type Bloom struct {
	bitset  Bitset `bson:"bitset"`
	sizeExp uint64 `bson:"sizeExp"`
	size    uint64 `bson:"size"`
	setLocs uint64 `bson:"setLocs"`
	shift   uint64 `bson:"shift"`
}

func (bl Bloom) Add(entry []byte) {
	l, h := bl.sipHash(entry)
	for i := uint64(0); i < bl.setLocs; i++ {
		bl.bitset.Set((h + i * l) & bl.size)
	}
}

func (bl Bloom) Has(entry []byte) bool {
	l, h := bl.sipHash(entry)
	for i := uint64(0); i < bl.setLocs; i++ {
		switch bl.bitset.IsSet((h + i * l) & bl.size) {
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
	bs    []uint64 `bson:"bs"`
	start uintptr  `bson:"start"`
}

func (bs *Bitset) Size(sz uint64) {
	(*bs).bs = make([]uint64, sz >> 6)
	(*bs).start = uintptr(unsafe.Pointer(&bs.bs[0]))
}

func (bs *Bitset) Set(idx uint64) {
	ptr := unsafe.Pointer(bs.start + uintptr(idx >> 3))
	*(*uint8)(ptr) |= mask[idx % 8]
}

func (bs *Bitset) Clear() {
	for i, _ := range (*bs).bs {
		(*bs).bs[i] = 0
	}
}

func (bs *Bitset) IsSet(idx uint64) bool {
	ptr := unsafe.Pointer(bs.start + uintptr(idx >> 3))
	r := ((*(*uint8)(ptr)) >> (idx % 8)) & 1
	return r == 1
}

func lineCounter(r io.ReadSeeker) (int, error) {
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}
		count += bytes.Count(buf[:c], lineSep)

		if c > 0 && buf[c - 1] != lineSep[0] {
			count++
		}

		if err == io.EOF {

			break
		}
	}
	_, err := r.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	return count, nil
}
func peek(buf *bytes.Buffer, b []byte) (int, error) {
	buf2 := bytes.NewBuffer(buf.Bytes())
	return buf2.Read(b)
}