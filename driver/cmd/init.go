// Copyright © 2018 munzli <manuel@monostream.com>
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
	"github.com/monostream/k8s-localflex-provisioner/driver/helper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the driver",
	Long: `Initialize the driver`,
	Run: func(cmd *cobra.Command, args []string) {

		cap := &helper.Capabilities{
			Attach: false,
		}
		res := helper.Response {
			Status:  helper.StatusSuccess,
			Message: "driver is available",
			Capabilities: cap,
		}

		helper.Handle(res)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
