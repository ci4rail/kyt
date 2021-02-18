/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	a "github.com/ci4rail/kyt/kyt-cli/internal/apply"
	e "github.com/ci4rail/kyt/kyt-cli/internal/errors"

	"github.com/spf13/cobra"
)

var (
	filename string
)

// almApplyCmd represents the almApply command
var almApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a resource by filename",
	Long: `Apply a configuration to a resource by filename. The resource name must be specified. This resource will be
created if it doesn't exist yet.

YAML format is accepted.`,
	Run: apply,
}

func apply(cmd *cobra.Command, args []string) {
	e.TokenConfigCheck()

	customerManifest := a.ReadCustomerManifest(filename)

	a.CustomerManifest(*customerManifest)
	fmt.Printf("Successfully deployed: %s\n", customerManifest.Application)
}

func init() {
	almCmd.AddCommand(almApplyCmd)

	almApplyCmd.Flags().StringVarP(&filename, "filename", "f", "", "that contains the configuration to apply (required)")
	err := almApplyCmd.MarkFlagRequired("filename")
	if err != nil {
		e.Er(err)
	}
}
