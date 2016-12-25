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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/casaplatform/casa/cmd/casa/cmd"
	"github.com/casaplatform/casa/cmd/casa/environment"

	// These imports are to enable these "plugins". This will change once
	// plugin support is released for Go.
	_ "github.com/casaplatform/daytime"
	_ "github.com/casaplatform/hue"
	_ "github.com/casaplatform/logger"
	_ "github.com/casaplatform/rules"
	_ "github.com/casaplatform/storage"
)

func main() {
	// Defer printing of the error, and avoid using log.Fatal or os.Exit
	// in order to allow any other deferred function to run.
	// See https://goo.gl/BLKYGv
	env := environment.Env
	env.WithOptions(environment.WithLogger(new(logLogger)))
	var err error
	defer func() {
		if err != nil {
			// Print out a stack trace
			if err, ok := err.(stackTracer); ok {
				env.Log("Casa has failed. Stack trace:")
				for _, f := range err.StackTrace() {
					fmt.Printf("%+s:%d", f, f)
				}
				os.Exit(1)
			} else {
				env.Log(err)
			}
		}
		// Do other stuff to shutdown
	}()
	err = cmd.RootCmd.Execute()
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// A simple casa.Logger that uses std lib log package
type logLogger struct{}

func (f logLogger) Log(a ...interface{}) {
	log.Println(a...)
}
