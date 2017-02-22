package bloomFilter

import (
	"testing"
	"bytes"
	"github.com/0x4139/bloomFilter"
)

func TestShouldCreateNewBloomFilter(t *testing.T) {
	_, err := bloomFilter.New(float64(1 << 16), float64(0.01))
	if err != nil {
		t.Error("Bloom filter could not be created", err.Error());
	}
}

func TestShouldCreateNewBloomFilterFromReader(t *testing.T) {
	value := "FOO"
	reader := bytes.NewReader([]byte(value))
	filter, err := bloomFilter.NewFromReader(reader, float64(0.01))
	if err != nil || filter == nil {
		t.Error("Bloom filter could not be created", err.Error());
	}
}

func TestShouldContainValueNewBloomFilterFromFile(t *testing.T) {
	filter, err := bloomFilter.NewFromFile("test_file", float64(0.01))
	if err != nil {
		t.Fatalf("Bloom filter could not be created", err.Error());
	}
	if !filter.Has([]byte("foo")) {
		t.Fatal("Bloom filter doesn't find the value!")
	}
}

func TestShouldContainValue(t *testing.T) {
	value := "fish"
	filter, _ := bloomFilter.New(float64(1 << 16), float64(0.01))
	filter.Add([]byte(value))
	if !filter.Has([]byte(value)) {
		t.Error("Bloom filter contains");
	}
}

func TestShouldNotContainValue(t *testing.T) {
	value := "fish"
	filter, _ := bloomFilter.New(float64(1 << 16), float64(0.01))
	filter.Add([]byte(value))
	if filter.Has([]byte("Fish")) {
		t.Error("Bloom filter does not contain " + value);
	}
}
