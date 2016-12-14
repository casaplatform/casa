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
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/casaplatform/casa"
	"github.com/casaplatform/casa/cmd/casa/environment"
	"github.com/casaplatform/mqtt"
	"github.com/gomqtt/broker"
	"github.com/gomqtt/packet"
	"github.com/gomqtt/transport"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the Casa server with an internal MQTT broker",
	RunE: func(cmd *cobra.Command, args []string) error {
		// why an "environment" package? http://www.jerf.org/iri/post/2929
		// Copy the global Environment so we get
		// the registered services. This can go away after plugins
		// are available in Go.

		// Use the environment packages's global Environment
		env := environment.Env
		brokerlogger := &brokerLogger{
			Logger: env.Logger,
		}

		// Create a new MessageBus by running our own MQTT broker
		bus, err := mqtt.New(
			mqtt.ListenOn(viper.GetString("MQTT.Listen")),
			mqtt.BrokerLogger(brokerlogger.Log),
		)

		if err != nil {
			return errors.Wrap(err, "Failed to create message bus")
		}

		// Set the remainder of the environment up
		env.WithOptions(
			environment.WithBus(bus),
			environment.WithViper(viper.GetViper()),
			environment.WithRegistry(environment.Env.ServiceRegistry),
			environment.WithLogger(new(logLogger)),
		)

		// Start listening for control-c and cleanly exit when called
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			sig := <-c
			env.Log("\nsignal: ", sig)
			var status int
			for key, s := range env.GetAllServices() {
				err := s.Stop()
				if err != nil {
					env.Log("Error stopping service", key, "::", err)
					status = 1
				}
			}

			err := env.MessageBus.Close()
			if err != nil {
				env.Log("Error closing bus", err)
				status = 1
			}
			os.Exit(status)
		}()

		for key, _ := range env.GetStringMap("Services") {
			if env.GetBool("Services." + key + ".Enabled") {
				config := env.Sub("Services." + key)
				config.Set("MQTT", env.Viper.Get("MQTT"))

				svc := env.GetService(key)
				if svc == nil {
					env.Log("Unsupported service: " + key)
					continue
				}

				svc.UseLogger(env)

				env.Log("Starting service: " + key)

				c := make(chan error, 1)
				go func(key string, ch chan error) {
					ch <- svc.Start(config)
				}(key, c)

				select {
				case err := <-c:
					if err != nil {
						env.Log("Failed starting service")
						continue
					}
					env.Log("Service started")
				case <-time.After(1 * time.Second):
					env.Log("Timeout while starting service")
				}
			}
		}

		for {
			// Loop forever!
			runtime.Gosched() // Play nice with go routines
		}
	},
}

// A simple casa.Logger that uses std lib log package
type logLogger struct{}

func (f logLogger) Log(a ...interface{}) {
	log.Println(a...)
}

// handles logging for the gomqqt.Broker
type brokerLogger struct {
	casa.Logger
}

// BrokerLogger logs errors from the MQTT broker
func (bl *brokerLogger) Log(event broker.LogEvent, client *broker.Client,
	packet packet.Packet, message *packet.Message, err error) {
	if err != nil {
		switch err.(type) {
		case transport.Error:
			e := err.(transport.Error)
			switch e.Code() {
			// NetworkError and ConnectionClose happen frequently
			// due to mobile clients dropping connections, etc.
			// I don't think we need to worry about them too much.
			case transport.NetworkError:
			case transport.ConnectionClose:

			default:
				bl.Logger.Log("Transport error", e.Code(), e.Err())
			}
		default:
			bl.Logger.Log(err)
		}
	}
}
