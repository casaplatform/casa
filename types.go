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

// Package casa defines interfaces and types for the casa home automation
// platform. These interfaces are subject to change without notice as casa is
// refined. Feel free to play with it now but wait until a v1 release before
// you depend on it for anything.
package casa

import "github.com/spf13/viper"

// Message represents an MQTT message
// Might be better to be an interface?
type Message struct {
	Topic   string
	Payload []byte
	Retain  bool
}

// MessageBus defines a system for passing Messages
type MessageBus interface {
	// Shutdown the MessageBus
	Close() error

	// NewClient returns a client that is connected to the MessageBus
	NewClient() MessageClient
}

// MessageClient sends and receives Messages to a MessageBus
type MessageClient interface {
	// Register a function to call when a message or error is received
	Handle(func(msg *Message, err error))

	// Publish a pre-composed message
	PublishMessage(message Message) error

	// Subscribes to the given topic
	Subscribe(topic string) error

	// Removes subscription for the topic
	Unsubscribe(topic string) error

	// Closes the client
	Close() error
}

// Service defines a type that  can subscribe to the message bus and perform
// actions.
type Service interface {
	// Start the service. Config is the config section for that service,
	// with the global MQTT config section appended.
	Start(config *viper.Viper) error

	// Specifies the logger to be used.
	UseLogger(logger Logger)

	// Stops the service.
	Stop() error
}

// Logger represents a way to log error messages for the user.
type Logger interface {
	Log(a ...interface{})
}
