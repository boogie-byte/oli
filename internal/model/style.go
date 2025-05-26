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

type style struct {
	// Base style
	base lipgloss.Style

	// Breadcrumbs
	breadcrumbs lipgloss.Style

	// Highlighted breadcrumb element
	breadcrumbHighlited lipgloss.Style

	// Done item
	itemDone lipgloss.Style

	// Bullet element
	bullet lipgloss.Style

	// Statusline (e.g. bottom bar info)
	statusline lipgloss.Style

	// Whitespace filing options for Place* methods
	whitespaceOpts []lipgloss.WhitespaceOption
}

type theme struct {
	background lipgloss.TerminalColor
	foreground lipgloss.TerminalColor
	grey       lipgloss.TerminalColor
	red        lipgloss.TerminalColor
	green      lipgloss.TerminalColor
	yellow     lipgloss.TerminalColor
	blue       lipgloss.TerminalColor
	purple     lipgloss.TerminalColor
	aqua       lipgloss.TerminalColor
	orange     lipgloss.TerminalColor
}

func (t *theme) style() style {
	base := lipgloss.NewStyle().
		Background(t.background).
		Foreground(t.foreground)

	return style{
		base: base,

		whitespaceOpts: []lipgloss.WhitespaceOption{
			lipgloss.WithWhitespaceBackground(t.background),
			lipgloss.WithWhitespaceForeground(t.foreground),
		},

		breadcrumbs: base.
			Foreground(t.grey).
			Italic(true).
			PaddingLeft(1),

		breadcrumbHighlited: base.
			Foreground(t.orange),

		itemDone: base.
			Foreground(t.grey).
			Strikethrough(true),

		bullet: base.
			Foreground(t.purple).
			Padding(0, 1),

		statusline: base.
			Background(t.grey).
			Padding(0, 1),
	}
}

var (
	// gruvbox-based default theme
	defaultTheme = theme{
		background: lipgloss.AdaptiveColor{
			Light: "#fbf1c7",
			Dark:  "#282828",
		},

		foreground: lipgloss.AdaptiveColor{
			Light: "#3c3836",
			Dark:  "#ebdbb2",
		},

		grey: lipgloss.AdaptiveColor{
			Light: "#928374",
			Dark:  "#928374",
		},

		// grey: lipgloss.AdaptiveColor{
		// 	Light: "#7c6f64",
		// 	Dark:  "#a89984",
		// },

		red: lipgloss.AdaptiveColor{
			Light: "#9d0006",
			Dark:  "#fb4934",
		},

		green: lipgloss.AdaptiveColor{
			Light: "#79740e",
			Dark:  "#b8bb26",
		},

		yellow: lipgloss.AdaptiveColor{
			Light: "#b57614",
			Dark:  "#fabd2f",
		},

		blue: lipgloss.AdaptiveColor{
			Light: "#076678",
			Dark:  "#83a598",
		},

		purple: lipgloss.AdaptiveColor{
			Light: "#8f3s71",
			Dark:  "#d3869b",
		},

		aqua: lipgloss.AdaptiveColor{
			Light: "#427b58",
			Dark:  "#8ec07c",
		},

		orange: lipgloss.AdaptiveColor{
			Light: "#af3a03",
			Dark:  "#fe8019",
		},
	}
)
