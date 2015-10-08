## bloomFilter: a bitset Bloom filter for go
===
This implements a fast bloom filter based on an 'unsafe' bitset

This uses [SipHash](https://en.wikipedia.org/wiki/SipHash) mostly for speed

###  Instalation

```sh
go get github.com/0x4139/bloomFilter
```

### Tests
Not many tests we're written :( sorry
```sh
go test
```
### Usage 
You can see the test example or 

```go
import "github.com/0x4139/bloomFilter"
func main{
    filter:=bloomFilter.New(float64(1<<16), float64(0.01)) //65535 items and 1% fail rate
    /* Other usages:
         New(float64(number_of_entries), float64(number_of_hashlocations))
         New(float64(100000), float64(2)) or New(float64(noentries), float64(nohashlocations))
         New(float64(100000), float64(0.05))
    */
}
```

### TODO
More tests
send pull requests please, love them