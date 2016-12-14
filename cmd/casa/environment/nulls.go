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

package environment

import (
	"os"

	"github.com/casaplatform/casa"
)

type nullSignalHandler struct{}

func (_ nullSignalHandler) HandleSignal(_ chan os.Signal) {}

type NullMessageBus struct{}

func (b NullMessageBus) Close() error                  { return nil }
func (b NullMessageBus) NewClient() casa.MessageClient { return nil }

type NullMessageStore struct{}

func (s NullMessageStore) Store(topic string, bus casa.MessageBus) error { return nil }
func (s NullMessageStore) Get(topic string) (casa.Message, error)        { return casa.Message{}, nil }
func (s NullMessageStore) Put(msg casa.Message) error                    { return nil }
func (s NullMessageStore) Close() error                                  { return nil }
