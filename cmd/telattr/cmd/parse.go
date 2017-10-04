package cmd

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/wenerme/telattr/genproto/v1/phonedata"
	"sort"
)

const (
	CMCC   byte = iota + 0x01 //中国移动
	CUCC                      //中国联通
	CTCC                      //中国电信
	CTCC_v                    //电信虚拟运营商
	CUCC_v                    //联通虚拟运营商
	CMCC_v                    //移动虚拟运营商
)

func vendor(b byte) string {
	switch b {
	case CMCC:
		return "中国移动"
	case CMCC_v:
		return "移动虚拟运营商"
	case CUCC:
		return "中国联通"
	case CUCC_v:
		return "联通虚拟运营商"
	case CTCC:
		return "中国电信"
	case CTCC_v:
		return "电信虚拟运营商"
	default:
		return fmt.Sprintf("未知运营商(%v)", b)
	}

}

func str(b []byte, i int) (s string, n int) {
	n = i
	//n = bytes.IndexByte(b[i:],0)
	for ; b[n] != 0; n++ {

	}

	s = string(b[i:n])
	n++
	return
}

func parseData(fn string) *phonedata.PhoneData {

	/*
	   | 4 bytes |                     <- phone.dat 版本号（如：1701即17年1月份）
	   ------------
	   | 4 bytes |                     <-  第一个索引的偏移
	   -----------------------
	   |  offset - 8            |      <-  记录区 - <省份>|<城市>|<邮编>|<长途区号>\0 - 山东|济南|250000|0531
	   -----------------------
	   |  index                 |      <-  索引区 - <手机号前七位><记录区的偏移><卡类型>
	   -----------------------
	*/

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	l := len(b)
	i := 4
	ind := int(binary.LittleEndian.Uint32(b[i:]))
	i += 4
	s := ""
	recs := make(map[int]*rec)

	for i < ind-8 {
		//end := bytes.Index(b[i:], []byte{0})
		//recs[i] = parseRec(string(b[i:end]))
		id := i
		s, i = str(b, i)
		recs[id] = parseRec(s)
	}

	//s, i = str(b, ind)
	_ = l

	idxs := make(map[int]*idx)
	idxes := make([]*idx, 0)
	for i = ind; i < l-8; i += 9 {
		idx := parseIdx(b[i:])
		idxs[i] = idx
		idxes = append(idxes, idx)
	}

	{
		data := &phonedata.PhoneData{}
		data.Version = string(b[:4])
		imap := make(map[int]int)
		for k, v := range recs {
			imap[k] = len(data.Records)
			data.Records = append(data.Records, &phonedata.Record{
				Province: v.prov,
				City:     v.city,
				Zip:      v.zip,
				Zone:     v.zone,
			})
		}

		for _, v := range idxs {
			data.Indexes = append(data.Indexes, &phonedata.Index{
				Prefix:      int32(v.prefix),
				RecordIndex: int32(imap[v.recIdx]),
				Vendor:      phonedata.Vendor(v.vendor),
			})
		}
		sort.Slice(data.Indexes, func(i, j int) bool { return data.Indexes[i].Prefix < data.Indexes[j].Prefix })

		return data
	}
}

func parseIdx(b []byte) *idx {
	return &idx{
		prefix: int(binary.LittleEndian.Uint32(b)),
		recIdx: int(binary.LittleEndian.Uint32(b[4:])),
		vendor: b[8],
	}
}

func parseRec(s string) *rec {
	// 山东|济南|250000|0531
	split := strings.Split(s, "|")
	return &rec{
		prov: split[0],
		city: split[1],
		zip:  split[2],
		zone: split[3],
	}
}

type idx struct {
	prefix int
	recIdx int
	rec    *rec
	vendor byte
}

func (self idx) String() string {
	return fmt.Sprintf("<%v><%v><%v>", self.prefix, self.recIdx, vendor(self.vendor))
}

type rec struct {
	zip  string
	zone string
	prov string
	city string
}

func (self rec) String() string {
	return fmt.Sprintf("%v|%v|%v|%v", self.prov, self.city, self.zip, self.zone)
}

func split(n int) []byte {
	var s []byte
	i := n
	for i != 0 {
		s = append(s, byte(i%10))
		i = i / 10
	}

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func build(recs map[int]*rec, idxs map[int]*idx) *tire {
	root := newTire()

	for _, idx := range idxs {
		b := split(idx.prefix)
		t := root.findOrCreate(b)
		t.idx = idx
		t.rec = recs[idx.recIdx]
	}

	return root
}

type tire struct {
	c        byte
	children []*tire
	idx      *idx
	rec      *rec
}
type tireOpt struct {
	create bool
	near   bool
}

func (self *tire) findOrCreate(b []byte) *tire {
	return self.findOpt(b, tireOpt{
		create: true,
	})
}
func (self *tire) findOpt(b []byte, opt tireOpt) *tire {
	if len(b) == 0 {
		return self
	}
	c := b[0]
	t := self.children[c]
	if t == nil {

		t = newTire()
		t.c = c
		self.children[c] = t
	}
	return t.findOrCreate(b[1:])
}
func newTire() *tire {
	return &tire{
		children: make([]*tire, 10, 10),
	}
}

type ByteSlice []byte

func (p ByteSlice) Len() int           { return len(p) }
func (p ByteSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p ByteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
