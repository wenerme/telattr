package telattr_test

import (
	"fmt"
	"github.com/wenerme/telattr"
	_ "github.com/wenerme/telattr/data_proto"
	"testing"
)

func BenchmarkFindPhone(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i = 0
		for p.Next() {
			i++
			_, err := telattr.Find(fmt.Sprintf("%s%d%s", "1897", i&10000, "45"))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
