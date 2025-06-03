// Copyright 2025 Sergey Vinogradov
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

package data

import "fmt"

type Status int

const (
	StatusNone Status = iota
	StatusToDo
	StatusDone
	StatusCanceled
	StatusWaiting
	StatusScheduled
)

func ParseStatus(s string) (Status, error) {
	switch s {
	case "NONE":
		return StatusNone, nil
	case "DONE":
		return StatusDone, nil
	case "CANC":
		return StatusCanceled, nil
	case "TODO":
		return StatusToDo, nil
	case "WAIT":
		return StatusWaiting, nil
	case "SCHD":
		return StatusScheduled, nil
	default:
		return -1, fmt.Errorf("unexpected status string: %s", s)
	}
}

func (s Status) String() string {
	switch s {
	case StatusNone:
		return "NONE"
	case StatusDone:
		return "DONE"
	case StatusCanceled:
		return "CANC"
	case StatusToDo:
		return "TODO"
	case StatusWaiting:
		return "WAIT"
	case StatusScheduled:
		return "SCHD"
	default:
		panic("unexpected status value")
	}
}
