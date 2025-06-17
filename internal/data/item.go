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

import (
	"encoding/xml"
	"io"
	"strconv"

	"github.com/google/uuid"
)

var (
	strTrue = strconv.FormatBool(true)
)

type Item struct {
	// Workspace the item belongs to
	workspace *Workspace

	// Item unique id
	id uuid.UUID

	// Parent item
	parent *Item

	// Linked list of child items
	head *Item
	tail *Item

	// Prevoius and next items in the siblings list
	prev *Item
	next *Item

	// data fields
	title     string
	status    Status
	collapsed bool
}

// Detach detaches the item from its parent and siblings.
func (i *Item) Detach() {
	if i.prev != nil {
		i.prev.next = i.next
	} else if i.parent != nil {
		i.parent.head = i.next
	}

	if i.next != nil {
		i.next.prev = i.prev
	} else if i.parent != nil {
		i.parent.tail = i.prev
	}

	i.parent = nil
	i.prev = nil
	i.next = nil
}

// MoveAbove moves item above the target.
func (i *Item) MoveAbove(target *Item) {
	i.Detach()

	i.parent = target.parent
	i.prev = target.prev
	i.next = target

	if target.prev != nil {
		target.prev.next = i
	} else {
		target.parent.head = i
	}

	target.prev = i
}

// MoveBelow moves item below the target.
func (i *Item) MoveBelow(target *Item) {
	i.Detach()

	i.parent = target.parent
	i.prev = target
	i.next = target.next

	if target.next != nil {
		target.next.prev = i
	} else {
		target.parent.tail = i
	}

	target.next = i
}

// Prepend places the provided item in the head position
// of the visitor's children list.
func (i *Item) Prepend(item *Item) {
	if i.head != nil {
		item.MoveAbove(i.head)
		return
	}

	item.Detach()

	item.parent = i
	i.head = item
	i.tail = item
}

// Append places the provided item in the tail position
// of the visitor's children list.
func (i *Item) Append(item *Item) {
	if i.tail != nil {
		item.MoveBelow(i.tail)
		return
	}

	item.Detach()

	item.parent = i
	i.head = item
	i.tail = item
}

// MoveUp places item before its previous sibling.
// If the prevoius sibling is nil, this method does nothing.
func (i *Item) MoveUp() {
	if i.prev != nil {
		i.MoveAbove(i.prev)
	}
}

// MoveDown places item after its next sibling.
// If the next sibling is nil, this method does nothing.
func (i *Item) MoveDown() {
	if i.next != nil {
		i.MoveBelow(i.next)
	}
}

func (i *Item) Demote() {
	prev := i.prev
	if prev == nil {
		return
	}

	prev.collapsed = false
	prev.Append(i)
}

func (i *Item) Promote() {
	// do not promote root item
	if i == i.workspace.root {
		return
	}

	// do not promote root children items
	if i.parent == i.workspace.root {
		return
	}

	i.MoveBelow(i.parent)
}

// PrevRow return the item on previous outline row, taking the
// "collapsed" state into consideration.
func (i *Item) PrevRow() *Item {
	prev := i.prev

	if prev == nil {
		parent := i.parent
		if parent == i.workspace.root {
			return nil
		}
		return parent
	}

	for {
		if prev.tail == nil || prev.collapsed {
			return prev
		}

		prev = prev.tail
	}
}

// NextRow return the item on next outline row, taking the
// "collapsed" state into consideration.
func (i *Item) NextRow() *Item {
	if head := i.head; head != nil && !i.collapsed {
		return head
	}

	if next := i.next; next != nil {
		return next
	}

	for p := i.parent; p != nil && p != i.workspace.root; p = p.parent {
		if p.next != nil {
			return p.next
		}
	}

	return nil
}

func (i *Item) Parent() *Item {
	return i.parent
}

func (i *Item) Head() *Item {
	return i.head
}

func (i *Item) Tail() *Item {
	return i.tail
}

// Title returns the item title.
func (i *Item) Title() string {
	return i.title
}

// Status returns the item status.
func (i *Item) Status() Status {
	return i.status
}

// Collapsed returns the item "collapsed" flag value.
func (i *Item) Collapsed() bool {
	return i.collapsed
}

// Depth returns the tree depth of the item relative to the
// workspace root. If the item is not in the workspace root,
// -1 is returned.
func (i *Item) Depth() int {
	depth := 0

	for {
		if i == i.workspace.root {
			break
		}

		if i.parent == nil {
			return -1
		}

		depth++

		i = i.parent
	}

	return depth
}

// ToDoStats returns the number of children in statuses "Done"
// or "Canceled" and the number of children in statuses other
// than "None".
func (item *Item) ToDoStats() (int, int) {
	var completed, total int
	for c := item.Head(); c != nil; c = c.Next() {
		s := c.Status()

		if s != StatusNone {
			total++
		}

		if s == StatusDone || s == StatusCanceled {
			completed++
		}
	}

	return completed, total
}

// DisplayedChildren returns a flattened list of non-collapsed
// child items.
func (i *Item) DisplayedChildren() []*Item {
	var items []*Item
	for c := i.Head(); c != nil; c = c.Next() {
		items = append(items, c)

		if !c.Collapsed() && c.Head() != nil {
			items = append(items, c.DisplayedChildren()...)
		}
	}
	return items
}

// RealRoot returns the root of the tree the item belongs to.
func (i *Item) RealRoot() *Item {
	r := i.workspace.root
	for {
		if r.parent == nil {
			return r
		}
		r = r.parent
	}
}

// Prev returns the previous sibling item.
func (i *Item) Prev() *Item {
	return i.prev
}

// Next returns the next sibling item.
func (i *Item) Next() *Item {
	return i.next
}

// SetTitle updates the item title value and marks the item as dirty.
func (i *Item) SetTitle(val string) {
	i.title = val
}

// SetStatus toggles the item "done" flag value and marks the
// item as dirty.
func (i *Item) SetStatus(s Status) {
	i.status = s
}

// SetCollapsed set the item "collapsed" flag value. If recursive is true
// it walks through the child items as well.
func (i *Item) SetCollapsed(value, recursive bool) {
	// collapse only items with children
	if value == true && i.head != nil {
		i.collapsed = true
	} else {
		i.collapsed = false
	}

	if recursive {
		for c := i.head; c != nil; c = c.next {
			c.SetCollapsed(value, recursive)
		}
	}
}

func newTrueAttr(name string) xml.Attr {
	return xml.Attr{
		Name:  xml.Name{Local: name},
		Value: strTrue,
	}
}

func (i *Item) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = xmlElemItem
	start.Attr = append(start.Attr, xml.Attr{
		Name:  xml.Name{Local: xmlItemAttrId},
		Value: i.id.String(),
	})

	if i.status != StatusNone {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: xmlItemAttrStatus},
			Value: i.status.String(),
		})
	}

	if i.collapsed {
		start.Attr = append(start.Attr, newTrueAttr(xmlItemAttrCollapsed))
	}

	if i == i.workspace.cursor {
		start.Attr = append(start.Attr, newTrueAttr(xmlItemAttrSelected))
	}

	if i == i.workspace.root {
		start.Attr = append(start.Attr, newTrueAttr(xmlItemAttrZoomedIn))
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeElement(i.title, xml.StartElement{Name: xml.Name{Local: xmlElemTitle}}); err != nil {
		return err
	}

	for c := i.head; c != nil; c = c.Next() {
		if err := e.Encode(c); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (i *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == xmlItemAttrId {
			var err error
			i.id, err = uuid.Parse(attr.Value)
			if err != nil {
				return err
			}
		}

		if attr.Name.Local == xmlItemAttrStatus {
			var err error
			i.status, err = ParseStatus(attr.Value)
			if err != nil {
				return err
			}
		}

		if attr.Name.Local == xmlItemAttrCollapsed {
			i.collapsed = true
		}

		if attr.Name.Local == xmlItemAttrSelected {
			i.workspace.cursor = i
		}

		if attr.Name.Local == xmlItemAttrZoomedIn {
			i.workspace.root = i
		}
	}

	for {
		tok, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch se := tok.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case xmlElemTitle:
				if err := d.DecodeElement(&i.title, &se); err != nil {
					return err
				}
			case xmlElemItem:
				c := i.workspace.NewItem("")
				if err := d.DecodeElement(c, &se); err != nil {
					return err
				}
				i.Append(c)
			default:
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if se.Name == start.Name {
				return nil
			}
		}
	}

	return nil
}
