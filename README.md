# telattr
Telephone Attribution

* Data source [xluohome/phonedata](https://github.com/xluohome/phonedata)
* Converted to protobuf
    * can use in web
    * no need to parse
* Use bind data, so, no need to keep the dat file
* Pre-generated protodata, csv
    * [Releases](https://github.com/wenerme/telattr/releases)

## Demo

```go
package main

import (
	"fmt"
	"github.com/wenerme/telattr"
	_ "github.com/wenerme/telattr/data_proto"
)

func main() {
	rec, err := telattr.Find("18957509123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", rec)
}
```

```
&telattr.Record{Version:"1707", Province:"浙江", City:"绍兴", Zip:"312000", Zone:"0575", VendorName:"中国电信"}
```


## Performance

A little bit faster than [xluohome/phonedata](https://github.com/xluohome/phonedata)

```
goos: darwin
goarch: amd64
pkg: github.com/xluohome/phonedata
BenchmarkFindPhone-8    10000000               181 ns/op

goos: darwin
goarch: amd64
pkg: github.com/wenerme/telattr
BenchmarkFindPhone-8    10000000               136 ns/op
```

## Dev
```bash
# Generate proto data
go-bindata -pkg data_proto -o data_proto/bindata.go phone.pb
```