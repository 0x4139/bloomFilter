package bloomFilter

import (
	"unsafe"
	"os"
	"bufio"
	"bytes"
	"io"
	"net/http"
	"time"
	"path/filepath"
	"github.com/jlaffaye/ftp"
)

const (
	ONE_IN_TEN_THOUSAND float64 = 0.01
	ONE_IN_ONE_HUNDRED_THOUSANDS float64 = 0.001
)

var CacheFolder string = "/tmp/bloomFilter/"
var TTL time.Duration = time.Hour * 24

// helper
var mask = []uint8{1, 2, 4, 8, 16, 32, 64, 128}

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
	filterSize, err := countLines(reader)
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

func NewFromUrl(url string, failRate float64) (*Bloom, error) {
	filename := newMd5FromString(url)
	//fileExtension := ".txt"
	tempFilePath := filepath.Join(CacheFolder, filename)
	fileInfo, err := os.Stat(tempFilePath)
	if err != nil || fileInfo.ModTime().Add(TTL).After(time.Now()) {
		err := os.MkdirAll(CacheFolder, 0755)
		if err != nil {
			return nil, err
		}
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			return nil, err
		}
		defer tempFile.Close()
		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return NewFromFile(tempFilePath, failRate)
}

func NewFromFTP(ftpAddress, username, password, ftpFilePath string, failRate float64) (filter *Bloom, err error) {
	filename := newMd5FromString(ftpFilePath)
	tempFilePath := filepath.Join(CacheFolder, filename)
	fileInfo, err := os.Stat(tempFilePath)
	if err != nil || fileInfo.ModTime().Add(TTL).After(time.Now()) {
		err := os.MkdirAll(CacheFolder, 0755)
		if err != nil {
			return
		}
		conn, err := ftp.Connect(ftpAddress)
		if err != nil {
			return
		}
		err = conn.Login(username, password)
		if err != nil {
			return
		}
		cd, err := conn.CurrentDir()
		if err != nil {
			return
		}
		ftpFile, err := conn.RetrFrom(cd + ftpFilePath, 0)
		if err != nil {
			return
		}
		out, err := os.Create(tempFilePath)
		if err != nil {
			return
		}
		_, err = io.Copy(out, ftpFile)
		if err != nil {
			return
		}
		out.Close()
		conn.Quit()
	}
	return NewFromFile(tempFilePath, failRate)
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

// TODO
func (bl Bloom) HasMd5(md5 string) bool {
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