package bloomFilter

import (
	"io"
	"bytes"
	"math"
	"crypto/md5"
	"encoding/hex"
)

func newMd5FromString(value string) string {
	md5Filename := md5.New()
	md5Filename.Write([]byte(value))
	return hex.EncodeToString(md5Filename.Sum(nil))
}

func countLines(r io.ReadSeeker) (int, error) {
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