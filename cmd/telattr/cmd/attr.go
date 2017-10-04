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

	"github.com/wenerme/telattr"
	_ "github.com/wenerme/telattr/data_proto"
	"fmt"
)

// attrCmd represents the attr command
var attrCmd = &cobra.Command{
	Use:     "attr number...",
	Aliases: []string{"a"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Query telephone attribution",
	Long:    `Query telephone attribution for give number`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			fmt.Printf("number:%v\n", v)
			rec, err := telattr.Find(v)
			if err != nil {
				fmt.Printf("error:%v\n", err.Error())
			} else {
				fmt.Printf(`vendor:%v
province:%v
city:%v
zone:%v
zip:%v
min:%v
max:%v
`, rec.VendorName, rec.Province, rec.City, rec.Zone, rec.Zip, rec.MinPrefix, rec.MaxPrefix)
				fmt.Println()
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(attrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// attrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// attrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
