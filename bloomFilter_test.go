package bloomFilter
import "testing"

func TestContains(t *testing.T) {
	filter := New(float64(1<<16), float64(0.01))
	filter.Add([]byte("fish"))
	if !filter.Has([]byte("fish")){
		t.Error("Bloom filter not working");
	}
}

func TestDoesNotContains(t *testing.T) {
	filter := New(float64(1<<16), float64(0.01))
	filter.Add([]byte("fish"))
	if filter.Has([]byte("Fish")){
		t.Error("Bloom filter not working");
	}
}
