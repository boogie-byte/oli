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

package model

import "github.com/charmbracelet/lipgloss"

const (
	black   = lipgloss.ANSIColor(0)
	red     = lipgloss.ANSIColor(1)
	green   = lipgloss.ANSIColor(2)
	yellow  = lipgloss.ANSIColor(3)
	blue    = lipgloss.ANSIColor(4)
	magenta = lipgloss.ANSIColor(5)
	cyan    = lipgloss.ANSIColor(6)
	white   = lipgloss.ANSIColor(7)
	grey    = lipgloss.ANSIColor(8)
)

var (
	styleBreadcrumbs = lipgloss.NewStyle().
				Foreground(grey).
				Italic(true).
				PaddingLeft(1)

	styleBreadcrumbHighlited = lipgloss.NewStyle().
					Foreground(magenta)

	styleItemDone = lipgloss.NewStyle().
			Foreground(grey).
			Strikethrough(true)

	styleStatusline = lipgloss.NewStyle().
			Reverse(true).
			Padding(0, 1)

	styleBullet = []lipgloss.Style{
		lipgloss.NewStyle().
			Foreground(yellow).
			Padding(0, 1),

		lipgloss.NewStyle().
			Foreground(green).
			Padding(0, 1),

		lipgloss.NewStyle().
			Foreground(cyan).
			Padding(0, 1),

		lipgloss.NewStyle().
			Foreground(blue).
			Padding(0, 1),

		lipgloss.NewStyle().
			Foreground(magenta).
			Padding(0, 1),

		lipgloss.NewStyle().
			Foreground(red).
			Padding(0, 1),
	}
)
