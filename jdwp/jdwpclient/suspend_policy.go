// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jdwpclient

import "fmt"

// SuspendPolicy describes what threads should be suspended on an event being
// raised.
type SuspendPolicy byte

const (
	// SuspendNone suspends no threads when a event is raised.
	SuspendNone = SuspendPolicy(0)
	// SuspendEventThread suspends only the event's thread when a event is raised.
	SuspendEventThread = SuspendPolicy(1)
	// SuspendAll suspends all threads when a event is raised.
	SuspendAll = SuspendPolicy(2)
)

func (s SuspendPolicy) String() string {
	switch s {
	case SuspendNone:
		return "SuspendNone"
	case SuspendEventThread:
		return "SuspendEventThread"
	case SuspendAll:
		return "SuspendAll"
	}
	return fmt.Sprint(int(s))
}
