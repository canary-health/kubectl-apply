// Copyright © 2018 Canary Health <sgraham@canaryhealth.com>
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
)

// deploymentCmd represents the deployment command
var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Applies `kind: Deployment` files to k8s cluster",
	Long:  "Applies `kind: Deployment` files using `apiVersion: apps/v1` to k8s cluster",
	Run:   cmdHandler.HandleDeployment,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)

	deploymentCmd.Flags().StringP("paramPath", "p", "", "Path to namespace params (e.g. /development/service)")

	deploymentCmd.Flags().StringP("image", "i", "", "Name of image to deploy")
}
