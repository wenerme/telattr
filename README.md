# telattr
Telephone Attribution

* Data source [xluohome/phonedata](https://github.com/xluohome/phonedata)
* Converted to protobuf
    * can use in web
    * no need to parse
* Use bind data, so, no need to keep the dat file
    * Use optimized pb data.
* Pre-generated protodata, csv
    * [Releases](https://github.com/wenerme/telattr/releases)
* Optimized index, remove continue prefix with same record index.
    * Original phone.dat 3.2M
    * PB phone.pb 3.8M
    * PB optimized phone-opt.pb 1.9M
* Useful command line tool

## CLI
```bash
go get -u github.com/wenerme/telattr/cmd/telattr
```

```bash
# Query
telattr a 1852159 15807212 15208231563
```
 
```
number:1852159
vendor:中国联通
province:上海
city:上海
zone:021
zip:200000
min:1852100
max:1852199

number:15807212
vendor:中国移动
province:湖北
city:荆州
zone:0716
zip:434000
min:1580721
max:1580721

number:15208231563
vendor:中国移动
province:四川
city:成都
zone:028
zip:610000
min:1520810
max:1520849
```

```bash
# Conversion
# Convert raw data to pb with optimize
telattr conv -f data phone.dat -t pb test.pb -o
# Conver pb to csv
# Output csv idx-test.csv, rec-test.csv
telattr conv -f pb phone.pb -t pb test.csv -o
```

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
&telattr.Record{Version:"1707", Province:"浙江", City:"绍兴", Zip:"312000", Zone:"0575", VendorName:"中国电信", MinPrefix:1895750, MaxPrefix:1895759}
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