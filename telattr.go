package telattr

import (
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	pb "github.com/wenerme/telattr/genproto/v1/phonedata"
	"sort"
	"strconv"
)

var ProtoData []byte
var PhoneData *pb.PhoneData

var ErrNotFound = errors.New("Not found")
var ErrInvalidNumber = errors.New("Invalid number")

func Find(num string) (*Record, error) {
	if len(num) < 7 {
		return nil, ErrInvalidNumber
	}
	phone_num, err := strconv.Atoi(num[:7])
	if err != nil {
		return nil, ErrInvalidNumber
	}
	phone_num_int32 := int32(phone_num)

	i := sort.Search(len(PhoneData.Indexes), func(i int) bool { return PhoneData.Indexes[i].Prefix >= phone_num_int32 })
	if i < 0 || i >= len(PhoneData.Indexes) {
		return nil, ErrNotFound
	}
	idx := PhoneData.Indexes[i]
	rec := PhoneData.Records[idx.RecordIndex]

	first, last := findRange(i)
	return &Record{
		Version:    PhoneData.Version,
		Province:   rec.Province,
		City:       rec.City,
		Zip:        rec.Zip,
		Zone:       rec.Zone,
		VendorName: VendorName(int(idx.Vendor)),

		MinPrefix: first.Prefix,
		MaxPrefix: last.Prefix,
	}, nil
}
func findRange(idx int) (*pb.Index, *pb.Index) {
	idxex := PhoneData.Indexes[idx]
	first := idxex
	last := idxex
	for i := idx; i >= 0; i -- {
		if idxex.RecordIndex == PhoneData.Indexes[i].RecordIndex {
			first = PhoneData.Indexes[i]
		} else {
			break
		}
	}
	for i := idx; i < len(PhoneData.Indexes); i ++ {
		if idxex.RecordIndex == PhoneData.Indexes[i].RecordIndex {
			last = PhoneData.Indexes[i]
		} else {
			break
		}
	}
	return first, last
}

func VendorName(v int) string {
	switch pb.Vendor(v) {
	case pb.Vendor_CMCC:
		return "中国移动"
	case pb.Vendor_CMCC_V:
		return "移动虚拟运营商"
	case pb.Vendor_CUCC:
		return "中国联通"
	case pb.Vendor_CUCC_V:
		return "联通虚拟运营商"
	case pb.Vendor_CTCC:
		return "中国电信"
	case pb.Vendor_CTCC_V:
		return "电信虚拟运营商"
	default:
		return fmt.Sprintf("未知运营商(%v)", v)
	}
}

type Record struct {
	Version string

	Province   string
	City       string
	Zip        string
	Zone       string
	VendorName string

	MinPrefix int32
	MaxPrefix int32
}

func MustInit() {
	err := doInit()
	if err != nil {
		panic(err)
	}
}
func doInit() error {
	if PhoneData == nil {
		if ProtoData == nil {
			return errors.New("Failed to init: no data")
		}
		data := &pb.PhoneData{}
		err := proto.Unmarshal(ProtoData, data)
		if err != nil {
			return err
		}
		PhoneData = data
	}
	return nil
}
