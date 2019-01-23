// Copyright © 2019 openSUSE opensuse-project@opensuse.org
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

	"github.com/kubic-project/caasp-init/pkg/certs"
	"github.com/kubic-project/caasp-init/pkg/config"
	"github.com/kubic-project/caasp-init/pkg/daemon"

	"github.com/spf13/cobra"
)

var (
	cfgFile              string
	registryConfigFolder = "/etc/docker"
)

const (
	longDescription = `Set the initial docker daemon mirrors configuration.

usage:

$ caasp-init

This will use default value for configuration file 

'/etc/kubic/kubic-init.yaml'

if you want to indicate a file run

$ caasp-init -c /path/to/config/file.yaml

If the configuration file has mirrors declared ot will generate the daemon.json
file with the following structure:

  {
    "registries": [
      {
      "Mirrors": [
        {
        "URL": "https://airgappedregistry.com"
        }
      ],
      "Prefix": "https://mycompany.registry.com"
      }
    ],
    "iptables":false,
    "log-level": "warn"
  }

If there is no mirror declared the configuration file will just be the default:

  {
    "iptables":false,
    "log-level": "warn"
    }

For help use 'caasp-init help'
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "caasp-init",
	Short: "Set the initial docker daemon mirrors configuration.",
	Long:  longDescription,
	RunE:  runE,
}

func runE(cmd *cobra.Command, args []string) error {
	err := os.MkdirAll(registryConfigFolder, 0644)
	if err != nil {
		return err
	}

	kubicConfig, err := config.FileAndDefaultsToKubicInitConfig(cfgFile)
	if err != nil {
		return err
	}

	err = daemon.WriteConfigFile(kubicConfig)
	if err != nil {
		return err
	}

	err = certs.WriteCertificates(kubicConfig)
	if err != nil {
		return err
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "/etc/kubic/kubic-init.yaml", "kubibc-init.yaml config file")
	rootCmd.AddCommand(newVersionCmd())
}
