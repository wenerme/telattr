// Copyright © 2017 wener <wenermail@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"context"
	"fmt"
	"github.com/wenerme/telattr/genproto/v1/phonedata"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"os"
	"encoding/csv"
	"github.com/wenerme/telattr"
)

var optimize = false
var from = ""
var to = ""

var data *phonedata.PhoneData

// convCmd represents the conv command
var convCmd = &cobra.Command{
	Use:   "conv [in] [out]",
	Args:  cobra.MinimumNArgs(2),
	Short: "Data format conversion",
	Long:  `Data format conversion between multi format`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(nil, "from", from)
		switch from {
		case "data":
			openDat(args[0], ctx)
		case "pb":
			openPb(args[0], ctx)
		default:
			fmt.Println("Invalid source format")
		}
		if optimize {
			optimizeIndex()
		}

		switch to {
		case "csv":
			writeCsv(args[1], ctx)
		case "pb":
			writePb(args[1], ctx)
		default:
			fmt.Println("Invalid dest format")
		}
	},
}

func optimizeIndex() {
	var last *phonedata.Index
	final := make([]*phonedata.Index, 0)
	final = append(final, data.Indexes[0])
	for _, v := range data.Indexes {
		if last != nil && (last.Prefix+1 != v.Prefix || last.RecordIndex != v.RecordIndex) {
			final = append(final, last, v)
		}
		last = v
	}
	final = append(final, last)
	data.Indexes = final
}

func openDat(fn string, ctx context.Context) {
	data = parseData(fn)
}
func openPb(fn string, ctx context.Context) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	data = &phonedata.PhoneData{}
	err = proto.Unmarshal(b, data)
	if err != nil {
		panic(err)
	}
}
func writePb(fn string, ctx context.Context) {
	b, err := proto.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fn, b, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
func writeCsv(fn string, ctx context.Context) {

	idx, err := os.OpenFile("idx-"+fn, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	rec, err := os.OpenFile("rec-"+fn, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	idxw := csv.NewWriter(idx)
	recw := csv.NewWriter(rec)

	for _, v := range data.Indexes {

		idxw.Write([]string{
			fmt.Sprint(v.Prefix),
			fmt.Sprint(v.RecordIndex),
			telattr.VendorName(int(v.Vendor)),
		})
	}

	for k, v := range data.Records {
		recw.Write([]string{
			fmt.Sprint(k),
			v.Province,
			v.City,
			v.Zip,
			v.Zone,
		})
	}
	idxw.Flush()
	recw.Flush()
}

func init() {
	RootCmd.AddCommand(convCmd)

	convCmd.Flags().BoolVarP(&optimize, "optimize", "o", false, "Optimize index")
	convCmd.Flags().StringVarP(&from, "from", "f", "", "From format")
	convCmd.Flags().StringVarP(&to, "to", "t", "", "To format")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
