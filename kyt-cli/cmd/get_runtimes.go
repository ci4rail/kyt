/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/ci4rail/kyt/kyt-cli/internal/alm"
	"github.com/ci4rail/kyt/kyt-cli/internal/auth"
	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"
	"github.com/spf13/cobra"
)

// runtimesCmd represents the runtimes command
var getRuntimesCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "Display alm-runtimes",
	Long: `Display alm-runtimes

Prints a table of the most important information of alm-runtimes`,
	Run: getRuntimes,
}

func getRuntimes(cmd *cobra.Command, args []string) {
	auth.TokenConfigCheck()

	runtimes := alm.FetchRuntimes(args)

	if len(runtimes) > 0 {
		switch o := output; o {
		case "json", "j":
			j, err := alm.ConvertToJSON(&runtimes)
			if err != nil {
				e.Er(err)
			}
			fmt.Println(j)
		case "yaml", "y":
			y, err := alm.ConvertToYaml(&runtimes)
			if err != nil {
				e.Er(err)
			}
			fmt.Println(y)
		case "wide", "w":
			// wide: Add here some more information for the table
			fmt.Printf("%-40s\t%-16s\n", "RUNTIME ID", "CONNECTION STATE")
			for _, runtime := range runtimes {
				fmt.Printf("%-40s\t%-16s\n", runtime.GetId(), runtime.GetNetwork())
			}
		case "short", "s":
			// short: only the most important information
			fmt.Printf("%-40s\t%-16s\n", "RUNTIME ID", "CONNECTION STATE")
			for _, runtime := range runtimes {
				fmt.Printf("%-40s\t%s\n", runtime.GetId(), runtime.GetNetwork())
			}
		default:
			fmt.Println("Error: Invalid output format given. See 'help' for more information.")
			os.Exit(1)
		}
	}
}

func init() {
	almGetCmd.AddCommand(getRuntimesCmd)
	getRuntimesCmd.Flags().StringVarP(&output, "output", "o", "short", "Output format. One of: short|s|json|j|yaml|y||wide|w")
}
