package bloomFilter

import (
	"testing"
	"bytes"
	"github.com/0x4139/bloomFilter"
)

func TestShouldCreateNew(t *testing.T) {
	_, err := bloomFilter.New(float64(1 << 16), bloomFilter.ONE_IN_TEN_THOUSAND)
	if err != nil {
		t.Error("Bloom filter could not be created", err.Error());
	}
}

func TestShouldCreateNewFromReadSeeker(t *testing.T) {
	value := "FOO"
	reader := bytes.NewReader([]byte(value))
	filter, err := bloomFilter.NewFromReadSeeker(reader, bloomFilter.ONE_IN_TEN_THOUSAND)
	if err != nil || filter == nil {
		t.Error("Bloom filter could not be created", err.Error());
	}
}

func TestShouldContainValueNewFromFile(t *testing.T) {
	filter, err := bloomFilter.NewFromFile("test_file", bloomFilter.ONE_IN_TEN_THOUSAND)
	if err != nil {
		t.Fatalf("Bloom filter could not be created: %s", err.Error());
	}
	if !filter.Has([]byte("foo")) {
		t.Fatal("Bloom filter doesn't find the value!")
	}
}

func TestShouldContainValueNewFromUrl(t *testing.T) {
	url := "http://tilinga:bilinga@monsterbox.nexthosting.ro/internal/suppress.txt"
	firstEmail := "001crouch@live.com"
	filter, err := bloomFilter.NewFromUrl(url, bloomFilter.ONE_IN_TEN_THOUSAND)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !filter.Has([]byte(firstEmail)) {
		t.Fatal("Bloom filter should containcontain " + firstEmail + ", but it doesn't ")
	}
}

func TestShouldContainValue(t *testing.T) {
	value := "fish"
	filter, _ := bloomFilter.New(float64(1 << 16), bloomFilter.ONE_IN_TEN_THOUSAND)
	filter.Add([]byte(value))
	if !filter.Has([]byte(value)) {
		t.Error("Bloom filter contains");
	}
}

func TestShouldNotContainValue(t *testing.T) {
	value := "fish"
	filter, _ := bloomFilter.New(float64(1 << 16), bloomFilter.ONE_IN_TEN_THOUSAND)
	filter.Add([]byte(value))
	if filter.Has([]byte("Fish")) {
		t.Error("Bloom filter does not contain " + value);
	}
}
