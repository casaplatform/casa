// Copyright © 2016 Casa Platform
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
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "casa",
	Short: "Casa is a home automation service backed by MQTT",
	Long: `Casa is a home automation service that uses MQTT for message
passing enabling compatibility with a wide range of other software and services.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.casa.toml first or ./.casa.toml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var configDefaults = map[string]interface{}{
	"Cores": 0,    // Use all cores
	"Debug": true, // Print debugging output

	"MQTT.Listen": "tcp://:1883",

	"Services.Logger.Enabled": true,
	"Services.Logger.Topics":  []string{"#"}, // Subscribe to all by default

	"Services.Hue.Enabled":  false,
	"Services.Hue.BridgeIP": "",
	"Services.Hue.User":     "",

	"Services.Storage.Enabled": "true",
	"Services.Storage.Topics":  []string{"#"},
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set defaults
	for key, value := range configDefaults {
		viper.SetDefault(key, value)
	}

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".casa") // name of config file (without extension)

	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AddConfigPath(".")     // optionally look for config in the working directory
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		//log.Fatal("Using config file:", viper.ConfigFileUsed())
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("No config file found, writing defaults to .casa.toml")

			f, err := os.Create("./.casa.toml")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			if err := toml.NewEncoder(w).Encode(viper.AllSettings()); err != nil {
				log.Fatal("Panic while encoding into TOML format.")
			}

			if err := w.Flush(); err != nil {
				log.Fatal("Error writing new config:", err)
			}

			if err := viper.ReadInConfig(); err != nil {
				log.Fatal("Error reading config:", err)
			}

			log.Println("Using new default config file:", viper.ConfigFileUsed())
		default:
			log.Fatal(err)
		}
	} else {
	}

	// Watch the config for changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Reloading config file:", viper.ConfigFileUsed())
			// Need to signal to other code that the config has changed
		}
	})

	// Limit the number of cores, if desired.
	if viper.GetInt("Cores") < 1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(viper.GetInt("Cores"))
	}
}
