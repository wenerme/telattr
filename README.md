# telattr
Telephone Attribution

* Data source [xluohome/phonedata](https://github.com/xluohome/phonedata)
* Converted to protobuf
    * can use in web
    * no need to parse
* Use bind data, so, no need to keep the dat file
* Pre-generated protodata, csv
    * [Releases](https://github.com/wenerme/telattr/releases)

## Dev
```bash
# Generate proto data
go-bindata -pkg data_proto -o data_proto/bindata.go phone.pb
```