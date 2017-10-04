package data_proto

import "github.com/wenerme/telattr"

// func Asset(name string) ([]byte, error) {
func init() {
	b, err := Asset("phone-opt-1707.pb")
	if err != nil {
		panic(err)
	}
	telattr.ProtoData = b
	telattr.MustInit()
}
