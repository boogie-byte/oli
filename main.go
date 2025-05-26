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

package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/boogie-byte/oli/internal/data"
	"github.com/boogie-byte/oli/internal/model"
)

func main() {
	directory := os.ExpandEnv("$HOME/.oli")
	if err := os.MkdirAll(directory, 0700); err != nil {
		log.Fatal(err)
	}

	w, err := data.LoadWorkspace(directory)
	if err != nil {
		log.Fatal(err)
	}

	m, err := model.NewOutline(w)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
