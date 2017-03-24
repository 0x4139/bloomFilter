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
	url := "http://tilinga:bilinga@monsterbox.nexthosting.ro/internal/sparta/19152_MDR-2139-life-university-2139_2017-02-15.txt"
	emailToCheck := "c3d8cbcbc9923178d93a0f53111edae0" // md5 entry
	filter, err := bloomFilter.NewFromUrl(url, bloomFilter.ONE_IN_TEN_THOUSAND)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !filter.Has([]byte(emailToCheck)) {
		t.Fatal("Bloom filter should containcontain " + emailToCheck + ", but it doesn't ")
	}
}

func TestShouldContainValueNewFromFtp(t *testing.T) {
	url := "ftp://68.64.169.66"
	username := "hitpath"
	password := "dUpn4mAk2vuFNtvF"
	ftpFilePath := "1127_FreshStartTax.zip"
	emailToCheck := "004ee@comcast.net"
	filter, err := bloomFilter.NewFromFTP(url, username, password, ftpFilePath, bloomFilter.ONE_IN_TEN_THOUSAND)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !filter.Has([]byte(emailToCheck)) {
		t.Fatal("Bloom filter should containcontain " + emailToCheck + ", but it doesn't ")
	}
}

func TestShouldContainValue(t *testing.T) {
	valueToCheck := "fish"
	filter, _ := bloomFilter.New(float64(1 << 16), bloomFilter.ONE_IN_TEN_THOUSAND)
	filter.Add([]byte(valueToCheck))
	if !filter.Has([]byte(valueToCheck)) {
		t.Error("Bloom filter contains");
	}
}

func TestShouldNotContainValue(t *testing.T) {
	valueToCheck := "fish"
	filter, _ := bloomFilter.New(float64(1 << 16), bloomFilter.ONE_IN_TEN_THOUSAND)
	filter.Add([]byte(valueToCheck))
	if filter.Has([]byte("Fish")) {
		t.Error("Bloom filter does not contain " + valueToCheck);
	}
}