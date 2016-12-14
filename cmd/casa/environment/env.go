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

// The environment package provides a way to setup a global environment and also
// have the ability to create local encironments for testing.
package environment

import (
	"os"
	"sync"

	"github.com/casaplatform/casa"
	"github.com/gomqtt/broker"
	"github.com/gomqtt/packet"
	"github.com/spf13/viper"
)

// Global Environment, use it if you want to or create your own with New()
var Env = New()

func RegisterService(name string, service casa.Service) {
	Env.AddService(name, service)
}

type Environment struct {
	*viper.Viper
	casa.MessageBus
	*ServiceRegistry
	SignalHandler
	BrokerLogger broker.Logger
	casa.Logger
}

type Option func(*Environment)

func (e *Environment) WithOptions(options ...Option) *Environment {

	for _, option := range options {
		option(e)
	}

	if e.MessageBus == nil {
		e.MessageBus = NullMessageBus{}
	}

	if e.SignalHandler == nil {
		e.SignalHandler = nullSignalHandler{}
	}

	if e.BrokerLogger == nil {
		e.BrokerLogger = nullBrokerLogger
	}

	if e.Logger == nil {
		e.Logger = nullLogger{}
	}

	if e.Viper == nil {
		e.Viper = viper.New()
	}

	// Might want to make a Null ServiceRegistry, not sure yet

	return e
}

func New(options ...Option) *Environment {
	e := &Environment{
		ServiceRegistry: &ServiceRegistry{
			services: make(map[string]casa.Service),
		},
	}

	return e.WithOptions(options...)
}

type nullLogger struct{}

func (_ nullLogger) Log(...interface{}) {}

func nullBrokerLogger(event broker.LogEvent, client *broker.Client,
	packet packet.Packet, message *packet.Message, err error) {
}
func WithLogger(l casa.Logger) Option {
	return func(e *Environment) {
		e.Logger = l
	}
}

func WithBrokerLogger(l broker.Logger) Option {
	return func(e *Environment) {
		e.BrokerLogger = l
	}
}

func WithRegistry(svc *ServiceRegistry) Option {
	return func(e *Environment) {
		e.ServiceRegistry = svc
	}
}

func WithBus(bus casa.MessageBus) Option {
	return func(e *Environment) {
		e.MessageBus = bus
	}
}

func WithHandler(handler SignalHandler) Option {
	return func(e *Environment) {
		e.SignalHandler = handler
	}
}

func WithViper(v *viper.Viper) Option {
	return func(e *Environment) {
		e.Viper = v
	}
}

type ServiceRegistry struct {
	sync.RWMutex
	services map[string]casa.Service
}

func (s *ServiceRegistry) AddService(name string, service casa.Service) {
	if s == nil {
		return
	}
	s.Lock()
	defer s.Unlock()
	s.services[name] = service
}

func (s *ServiceRegistry) RemoveService(name string) {
	if s == nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	delete(s.services, name)
}

func (s *ServiceRegistry) GetAllServices() map[string]casa.Service {
	if s == nil {
		return nil
	}

	s.RLock()
	defer s.RUnlock()
	return s.services
}

func (s *ServiceRegistry) GetService(name string) casa.Service {
	if s == nil {
		return nil
	}

	s.RLock()
	defer s.RUnlock()
	return s.services[name]
}

type SignalHandler interface {
	HandleSignal(chan os.Signal)
}
