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

type Outline struct {
	workspace *data.Workspace

	windowWidth  int
	windowHeight int

	textInput textinput.Model

	commandMode    commandMode
	itemMode       itemMode
	itemStatusMode itemStatusMode

	statusLine string
}

func NewOutline(workspace *data.Workspace) (*Outline, error) {
	m := &Outline{
		workspace: workspace,
	}

	m.textInput = textinput.New()
	m.textInput.SetValue(workspace.Cursor().Title())
	m.textInput.Prompt = ""
	m.textInput.Focus()

	m.commandMode = commandMode{m}
	m.itemMode = itemMode{m}
	m.itemStatusMode = itemStatusMode{m}

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

func getStatus(item *data.Item) string {
	if s := item.Status(); s != data.StatusNone {
		return styleItemStatus[s].Render(s.String())
	}

	return ""
}

func getItemStyle(item *data.Item) lipgloss.Style {
	switch item.Status() {
	case data.StatusDone, data.StatusCanceled:
		return styleItemComplete
	default:
		return styleItemNormal
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

func (m *Outline) moveCursor(item *data.Item) (tea.Model, tea.Cmd) {
	if item == nil {
		return m, nil
	}

	m.saveCurrentTitle()
	m.updateTextInput(item)
	m.textInput.CursorEnd()

	m.workspace.SetCursor(item)

	return m, nil
}

func (m *Outline) cursorUp() (tea.Model, tea.Cmd) {
	item := m.workspace.Cursor().PrevRow()
	return m.moveCursor(item)
}

func (m *Outline) cursorDown() (tea.Model, tea.Cmd) {
	item := m.workspace.Cursor().NextRow()
	return m.moveCursor(item)
}

func (m *Outline) cursorToParent() (tea.Model, tea.Cmd) {
	parent := m.workspace.Cursor().Parent()
	if parent != m.workspace.Root() {
		return m.moveCursor(parent)
	}

	return m, nil
}

func (m *Outline) cursorToTail() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()

	tail := cur.Tail()
	if tail != nil {
		cur.SetCollapsed(false, false)
		return m.moveCursor(tail)
	}

	return m, nil
}

func (m *Outline) zoomIn() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	if cur.Head() == nil {
		return m, nil
	}

	m.workspace.SetRoot(cur)
	m.moveCursor(cur.Head())

	return m, nil
}

func (m *Outline) zoomOut() (tea.Model, tea.Cmd) {
	root := m.workspace.Root()
	if root.Parent() == nil {
		return m, nil
	}

	m.workspace.SetRoot(root.Parent())

	if root.Collapsed() {
		m.moveCursor(root)
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

func (m *Outline) toggleItemFolded(recursive bool) (tea.Model, tea.Cmd) {
	collapsed := m.workspace.Cursor().Collapsed()
	m.workspace.Cursor().SetCollapsed(!collapsed, recursive)

	return m, nil
}

func (m *Outline) toggleBranchCollapsed() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Outline) toggleRowDone() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	if cur.Status() == data.StatusDone {
		cur.SetStatus(data.StatusNone)
	} else {
		cur.SetStatus(data.StatusDone)
	}

	return m.moveCursor(cur.Next())
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

func (m *Outline) deleteItem(recursive bool) (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()

	nextSelected := cur.Next()
	if nextSelected == nil {
		nextSelected = cur.Prev()
	}

	if nextSelected == nil {
		nextSelected = cur.Parent()
	}

	if nextSelected == m.workspace.Root() {
		return m, nil
	}

	if cur.Head() != nil && !recursive {
		m.statusLine = styleStatusLineError.Render("Item has children, use C-c D for recursive deletion")
		return m, nil
	}

	cur.Detach()

	return m.moveCursor(nextSelected)
}

func (m *Outline) addSibling() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	next := m.workspace.NewItem("")

	if cur.Status() != data.StatusNone {
		next.SetStatus(data.StatusToDo)
	}

	next.MoveBelow(cur)

	return m.moveCursor(next)
}

func (m *Outline) addChild() (tea.Model, tea.Cmd) {
	cur := m.workspace.Cursor()
	next := m.workspace.NewItem("")
	tail := cur.Tail()

	if tail != nil && tail.Status() != data.StatusNone {
		next.SetStatus(data.StatusToDo)
	}

	cur.SetCollapsed(false, false)
	cur.Append(next)

	return m.moveCursor(next)
}

func (m *Outline) save() (tea.Model, tea.Cmd) {
	m.saveCurrentTitle()

	err := m.workspace.Save()
	if err != nil {
		m.statusLine = styleStatusLineError.Render(err.Error())
	} else {
		m.statusLine = styleStatusLineMessage.Render("Saved!")
	}

	return m, nil
}

func (m *Outline) resetStatusLineMessage() (tea.Model, tea.Cmd) {
	m.statusLine = ""
	return m, nil
}

func (m *Outline) updateWindowSize(msg tea.WindowSizeMsg) {
	m.windowWidth = msg.Width
	m.windowHeight = msg.Height
}

func (m *Outline) Init() tea.Cmd {
	return nil
}

func (m *Outline) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSize(msg)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlX:
			m.statusLine = m.commandMode.statusLine()
			return m.commandMode, nil
		case tea.KeyCtrlC:
			m.statusLine = m.itemMode.statusLine()
			return m.itemMode, nil
		case tea.KeyEsc:
			return m.resetStatusLineMessage()
		case tea.KeyCtrlUp:
			return m.cursorUp()
		case tea.KeyCtrlDown:
			return m.cursorDown()
		case tea.KeyCtrlLeft:
			return m.cursorToParent()
		case tea.KeyCtrlRight:
			return m.cursorToTail()
		case tea.KeyCtrlShiftUp:
			return m.moveRowUp()
		case tea.KeyCtrlShiftDown:
			return m.moveRowDown()
		case tea.KeyCtrlShiftRight:
			return m.demoteRow()
		case tea.KeyCtrlShiftLeft:
			return m.promoteRow()
		case tea.KeyTab:
			return m.addSibling()
		case tea.KeyShiftTab:
			return m.addChild()
		default:
			return m.updateRow(message)
		}
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

func (m *Outline) renderItemEntry(item *data.Item) string {
	bullet := getBullet(item)
	bullet = styleBullet[(item.Depth()-1)%len(styleBullet)].Render(bullet)

	status := getStatus(item)

	padding := getLinePadding(item)

	var title string
	if m.workspace.Cursor() == item {
		m.textInput.TextStyle = getItemStyle(item)
		title = m.textInput.View()
	} else {
		title = item.Title()

		maxTitleWidth := m.getMaxTitleWidth(padding)
		title = runewidth.Truncate(title, maxTitleWidth, "...")
		title = getItemStyle(item).Render(title)
	}

	var todoStats string
	if completed, total := item.ToDoStats(); completed != 0 || total != 0 {
		todoStats = fmt.Sprintf("(%d/%d)", completed, total)
		todoStats = styleTodoStats.Render(todoStats)
	}

	itemRow := lipgloss.JoinHorizontal(lipgloss.Top, bullet, status, title, todoStats)
	itemRow = lipgloss.PlaceHorizontal(
		m.windowWidth-padding,
		lipgloss.Left,
		itemRow,
	)

	return itemRow
}

func (m *Outline) renderItemList() string {
	var itemEntries []string
	for _, item := range m.workspace.Root().DisplayedChildren() {
		itemEntry := m.renderItemEntry(item)
		itemEntries = append(itemEntries, itemEntry)
	}

	items := lipgloss.JoinVertical(lipgloss.Right, itemEntries...)
	items = lipgloss.PlaceVertical(
		m.windowHeight-4,
		lipgloss.Top,
		items,
	)

	return items
}

func (m *Outline) renderStatusLine() string {
	return lipgloss.PlaceHorizontal(m.windowWidth, lipgloss.Top, m.statusLine)
}

func (m *Outline) View() string {
	// Wait for the window size to be set
	if m.windowWidth == 0 || m.windowHeight == 0 {
		return ""
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderBreadcrumbs(),
		m.renderItemList(),
		m.renderStatusLine(),
	)
}

type commandMode struct {
	*Outline
}

func (commandMode) statusLine() string {
	return "command: [q]uit without saving  [s]ave file"
}

func (m commandMode) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSize(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Outline.statusLine = ""
			return m.Outline, nil
		case "q":
			m.Outline.statusLine = ""
			return m.Outline, tea.Quit
		case "s":
			m.Outline.statusLine = ""
			m.save()
		default:
			return m, nil
		}
	}

	return m.Outline, nil
}

type itemMode struct {
	*Outline
}

func (itemMode) statusLine() string {
	return "item: [d]elete  [D]elete recursive  [f]old  [F]old recursive  change [s]tatus  [z]oom in  [Z]oom out"
}

func (m itemMode) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSize(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Outline.statusLine = ""
			return m.Outline, nil
		case "d":
			return m.deleteItem(false)
		case "D":
			return m.deleteItem(true)
		case "f":
			return m.toggleItemFolded(false)
		case "F":
			return m.toggleItemFolded(true)
		case "s":
			m.Outline.statusLine = m.Outline.itemStatusMode.statusLine()
			return m.Outline.itemStatusMode, nil
		case "z":
			m.Outline.statusLine = ""
			m.zoomIn()
		case "Z":
			m.Outline.statusLine = ""
			m.zoomOut()
		default:
			return m, nil
		}
	}

	return m.Outline, nil
}

type itemStatusMode struct {
	*Outline
}

func (itemStatusMode) statusLine() string {
	return "item status: [n]one  [t]odo  [d]one  [c]canceled  [w]aiting  [s]cheduled"
}

func (m itemStatusMode) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSize(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Outline.statusLine = ""
			return m.Outline, nil
		case "n":
			m.Outline.statusLine = ""
			m.Outline.workspace.Cursor().SetStatus(data.StatusNone)
		case "t":
			m.Outline.statusLine = ""
			m.Outline.workspace.Cursor().SetStatus(data.StatusToDo)
		case "d":
			m.Outline.statusLine = ""
			m.Outline.workspace.Cursor().SetStatus(data.StatusDone)
		case "c":
			m.Outline.statusLine = ""
			m.Outline.workspace.Cursor().SetStatus(data.StatusCanceled)
		case "w":
			m.Outline.statusLine = ""
			m.Outline.workspace.Cursor().SetStatus(data.StatusWaiting)
		case "s":
			m.Outline.statusLine = ""
			m.Outline.workspace.Cursor().SetStatus(data.StatusScheduled)
		default:
			return m, nil
		}
	}

	return m.Outline, nil
}
