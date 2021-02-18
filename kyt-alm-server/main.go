/*
Copyright © 2021 Ci4Rail GmbH

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

package main

import (
	"log"
	"os"

	"github.com/ci4rail/kyt/kyt-alm-server/cmd"
)

const (
	envIotHubConnectionsString = "IOTHUB_SERVICE_CONNECTION_STRING"
)

func main() {
	versionArgFound := false
	for _, v := range os.Args {
		if v == "version" || v == "help" || v == "--help" || v == "-h" {
			versionArgFound = true
		}
	}
	if !versionArgFound {
		_, ok := os.LookupEnv(envIotHubConnectionsString)

		if !ok {
			log.Fatalf("Error: environment variable %s missing", envIotHubConnectionsString)
		}
	}
	cmd.Execute()
}
