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

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/boogie-byte/oli/internal/data"
)

const (
	bulletFilledCircle  = "●" // U+25CF
	bulletEmptyCircle   = "◯" // U+25EF
	bulledTriangleRight = "▶" // U+25B6
	bulletTriangleDown  = "▼" // U+25BC

	prefixWitdh = 3
)

var (
	statusLineResetTimeout = 3 * time.Second
)

type statusLineResetMessage struct{}

func tickStatusLineReset() tea.Cmd {
	return tea.Tick(statusLineResetTimeout, func(_ time.Time) tea.Msg {
		return statusLineResetMessage{}
	})
}

type Outline struct {
	workspace *data.Workspace

	windowWidth  int
	windowHeight int

	commandMode bool
	textInput   textinput.Model

	errorMessage   string
	successMessage string
	infoMessage    string
}

func NewOutline(workspace *data.Workspace) (*Outline, error) {
	m := &Outline{
		workspace: workspace,
	}

	m.textInput = textinput.New()
	m.textInput.SetValue(workspace.Cursor().Title())
	m.textInput.Prompt = ""
	m.textInput.Focus()

	return m, nil
}

func getLinePadding(n *data.Item) int {
	return 2 * n.Depth()
}

func getBullet(item *data.Item) string {
	switch {
	case item.Head() == nil:
		return bulletFilledCircle
	case item.Collapsed():
		return bulledTriangleRight
	default:
		return bulletTriangleDown
	}
}

func (m *Outline) getMaxTitleWidth(padding int) int {
	return m.windowWidth - padding - prefixWitdh
}

func (m *Outline) breadcrumbs() string {
	var breadcrumbs string
	for p := m.workspace.Root().Parent(); p != nil; p = p.Parent() {
		breadcrumbs = p.Title() + " / " + breadcrumbs
	}

	return breadcrumbs
}

func (m *Outline) statusLine() string {
	switch {
	case m.errorMessage != "":
		return m.errorMessage
	case m.successMessage != "":
		return m.successMessage
	case m.infoMessage != "":
		return m.infoMessage
	default:
		return fmt.Sprintf(
			"%p <- %p -> %p",
			m.workspace.Cursor().Prev(),
			m.workspace.Cursor(),
			m.workspace.Cursor().Next(),
		)
	}
}

// Movement

func (m *Outline) saveCurrentTitle() {
	m.workspace.Cursor().SetTitle(m.textInput.Value())
}

func (m *Outline) updateTextInput(n *data.Item) {
	padding := getLinePadding(n)
	maxWidth := m.getMaxTitleWidth(padding)

	m.textInput.Width = 0
	if runewidth.StringWidth(n.Title()) > maxWidth {
		m.textInput.Width = maxWidth - 1 // -1 to show cursor
	}
	m.textInput.SetValue(n.Title())
}

func (m *Outline) moveCursor(item *data.Item, pos int) (tea.Model, tea.Cmd) {
	if item == nil {
		return m, nil
	}

	m.saveCurrentTitle()
	m.updateTextInput(item)
	if pos < 0 {
		m.textInput.CursorEnd()
	} else {
		m.textInput.SetCursor(pos)
	}

	m.workspace.SetCursor(item)

	return m, nil
}

func (m *Outline) cursorUp(pos int) (tea.Model, tea.Cmd) {
	item := m.workspace.Cursor().PrevRow()
	return m.moveCursor(item, pos)
}

func (m *Outline) cursorDown(pos int) (tea.Model, tea.Cmd) {
	item := m.workspace.Cursor().NextRow()
	return m.moveCursor(item, pos)
}

func (m *Outline) cursorHead() (tea.Model, tea.Cmd) {
	item := m.workspace.Cursor().Parent().Head()
	return m.moveCursor(item, -1)
}

func (m *Outline) cursorTail() (tea.Model, tea.Cmd) {
	item := m.workspace.Cursor().Parent().Tail()
	return m.moveCursor(item, -1)
}

func (m *Outline) zoomIn() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	if cur.Head() == nil {
		return m, nil
	}

	m.workspace.SetRoot(cur)
	m.moveCursor(cur.Head(), -1)

	return m, nil
}

func (m *Outline) zoomOut() (tea.Model, tea.Cmd) {
	root := m.workspace.Root()
	if root.Parent() == nil {
		return m, nil
	}

	m.workspace.SetRoot(root.Parent())

	if root.Collapsed() {
		m.moveCursor(root, -1)
	}

	return m, nil
}

// Row organizing

func (m *Outline) moveRowUp() (tea.Model, tea.Cmd) {
	m.workspace.Cursor().MoveUp()

	return m, nil
}

func (m *Outline) moveRowDown() (tea.Model, tea.Cmd) {
	m.workspace.Cursor().MoveDown()

	return m, nil
}

func (m *Outline) moveRowToHead() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	cur.Parent().Prepend(cur)

	return m, nil
}

func (m *Outline) moveRowToTail() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	cur.Parent().Append(cur)

	return m, nil
}

func (m *Outline) toggleRowCollapsed(recursive bool) (tea.Model, tea.Cmd) {
	collapsed := m.workspace.Cursor().Collapsed()
	m.workspace.Cursor().SetCollapsed(!collapsed, recursive)

	return m, nil
}

func (m *Outline) toggleBranchCollapsed() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Outline) toggleRowDone() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	cur.ToggleDone()

	return m.moveCursor(cur.Next(), -1)
}

func (m *Outline) demoteRow() (tea.Model, tea.Cmd) {
	m.saveCurrentTitle()

	cur := m.workspace.Cursor()
	cur.Demote()

	m.updateTextInput(cur)

	return m, nil
}

func (m *Outline) promoteRow() (tea.Model, tea.Cmd) {
	m.saveCurrentTitle()

	cur := m.workspace.Cursor()
	cur.Promote()

	m.updateTextInput(cur)

	return m, nil
}

func (m *Outline) updateRow(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *Outline) deleteEmptyRow(message tea.Msg) (tea.Model, tea.Cmd) {
	if m.textInput.Position() != 0 {
		return m.updateRow(message)
	}

	prevRow := m.workspace.Cursor().PrevRow()
	if prevRow == nil {
		return m, nil
	}

	prevRowPos := len(prevRow.Title())

	if v := m.textInput.Value(); v != "" {
		prevRow.SetTitle(prevRow.Title() + v)
	}

	// promote children if they could
	// be attached to previous row
	cur := m.workspace.Cursor()
	for c := cur.Head(); c != nil; c = c.Next() {
		prevRow.Append(c)
	}

	model, cmd := m.moveCursor(prevRow, prevRowPos)

	cur.Detach()

	return model, cmd
}

func (m *Outline) addRow() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	val := m.textInput.Value()

	if val == "" {
		if cur.Depth() > 1 {
			cur.MoveBelow(cur.Parent())
		}
		return m, nil
	}

	pos := m.textInput.Position()

	cur.SetTitle(val[:pos])
	m.updateTextInput(cur)

	next := m.workspace.NewItem(val[pos:])
	next.MoveBelow(cur)

	return m.moveCursor(next, -1)
}

func (m *Outline) save(quitAfterSave bool) (tea.Model, tea.Cmd) {
	m.saveCurrentTitle()

	err := m.workspace.Save()
	if err != nil {
		m.errorMessage = err.Error()
	} else if quitAfterSave {
		return m, tea.Quit
	} else {
		m.successMessage = "Saved!"
	}
	return m, tickStatusLineReset()
}

func (m *Outline) resetStatusLineMessage() (tea.Model, tea.Cmd) {
	m.errorMessage = ""
	m.successMessage = ""
	m.infoMessage = ""
	return m, nil
}

func (m *Outline) enterCommandMode() (tea.Model, tea.Cmd) {
	m.commandMode = true
	return m, nil
}

func (m *Outline) command(s string) (tea.Model, tea.Cmd) {
	m.commandMode = false

	switch s {
	// Misc
	case "q":
		return m.save(true)
	case "s":
		return m.save(false)
	// Zoom
	case "i":
		return m.zoomIn()
	case "o":
		return m.zoomOut()
	// Adding and modifying
	case "x":
		return m.toggleRowDone()
	}

	return m, nil
}

func (m *Outline) Init() tea.Cmd {
	return nil
}

func (m *Outline) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tea.KeyMsg:
		if m.commandMode {
			return m.command(msg.String())
		}

		switch msg.Type {
		case tea.KeyCtrlX:
			return m.enterCommandMode()
		case tea.KeyUp:
			return m.cursorUp(-1)
		case tea.KeyDown:
			return m.cursorDown(-1)
		case tea.KeyShiftUp:
			return m.cursorHead()
		case tea.KeyShiftDown:
			return m.cursorTail()
		case tea.KeyCtrlUp:
			return m.moveRowUp()
		case tea.KeyCtrlDown:
			return m.moveRowDown()
		case tea.KeyCtrlShiftUp:
			return m.moveRowToHead()
		case tea.KeyCtrlShiftDown:
			return m.moveRowToTail()
		case tea.KeyCtrlRight:
			return m.demoteRow()
		case tea.KeyCtrlLeft:
			return m.promoteRow()
		case tea.KeyBackspace:
			return m.deleteEmptyRow(message)
		case tea.KeyEnter:
			return m.addRow()
		case tea.KeyTab:
			return m.toggleRowCollapsed(false)
		case tea.KeyShiftTab:
			return m.toggleRowCollapsed(true)
		default:
			return m.updateRow(message)
		}
	case statusLineResetMessage:
		return m.resetStatusLineMessage()
	}

	return m, nil
}

func (m *Outline) renderBreadcrumbs() string {
	breadcrumbs := lipgloss.JoinHorizontal(
		lipgloss.Top,
		styleBreadcrumbs.Render(m.breadcrumbs()),
		styleBreadcrumbHighlited.Render(m.workspace.Root().Title()),
	)

	breadcrumbs = runewidth.Truncate(breadcrumbs, m.windowWidth-2, "...")

	breadcrumbs = lipgloss.PlaceHorizontal(
		m.windowWidth,
		lipgloss.Left,
		breadcrumbs,
	)

	breadcrumbs = lipgloss.PlaceVertical(
		3,
		lipgloss.Center,
		breadcrumbs,
	)

	return breadcrumbs
}

func (m *Outline) renderItems() string {
	var itemStrs []string
	for _, item := range m.workspace.Root().DisplayedChildren() {
		padding := getLinePadding(item)

		var title string
		if m.workspace.Cursor() == item {
			if item.Done() {
				m.textInput.TextStyle = styleItemDone
			} else {
				m.textInput.TextStyle = lipgloss.NewStyle()
			}
			title = m.textInput.View()
		} else {
			title = item.Title()

			maxTitleWidth := m.getMaxTitleWidth(padding)
			title = runewidth.Truncate(title, maxTitleWidth, "...")

			if item.Done() {
				title = styleItemDone.Render(title)
			}
		}

		bullet := getBullet(item)
		bullet = styleBullet[(item.Depth()-1)%len(styleBullet)].Render(bullet)

		itemRow := lipgloss.JoinHorizontal(lipgloss.Top, bullet, title)
		itemRow = lipgloss.PlaceHorizontal(
			m.windowWidth-padding,
			lipgloss.Left,
			itemRow,
		)

		itemStrs = append(itemStrs, itemRow)
	}

	items := lipgloss.JoinVertical(lipgloss.Right, itemStrs...)
	items = lipgloss.PlaceVertical(
		m.windowHeight-4,
		lipgloss.Top,
		items,
	)

	return items
}

func (m *Outline) renderStatusLine() string {
	return styleStatusline.
		Width(m.windowWidth).
		Render(m.statusLine())
}

func (m *Outline) View() string {
	// Wait for the window size to be set
	if m.windowWidth == 0 || m.windowHeight == 0 {
		return ""
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderBreadcrumbs(),
		m.renderItems(),
		m.renderStatusLine(),
	)
}
