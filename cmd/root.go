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
	"fmt"
	"os"

	"github.com/canary-health/kubectl-apply/config"
	"github.com/canary-health/kubectl-apply/handler"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var file string
var labels string
var namespace string
var kubeconfig string

// var c config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-apply",
	Short: "Create or update resources, similar to `kubectl apply`.",
	Long:  "Create or update resources, similar to `kubectl apply`. Implementing k8s client-go. Limited supported resources kinds. See docs for more.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var c *config.Config

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-apply.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "kubeconfig file (default is $KUBECONFIG set in build container or $HOME/.kube/config)")
	// rootCmd.MarkFlagRequired("kubeconfig")

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "absolute path to deployment file")
	rootCmd.MarkFlagRequired("file")

	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "development", "kubernetes cluster namespace")

	rootCmd.PersistentFlags().StringVarP(&labels, "labels", "l", "", "metadata.labels string comma delimited key:values (e.g. key1:value1,key2:value2)")
}

var cmdHandler = handler.NewHandler(c)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kubectl-apply" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kubectl-apply")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
