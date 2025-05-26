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

package data_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/boogie-byte/oli/internal/data"
)

func TestNewWorkspace(t *testing.T) {
	w := data.NewWorkspace("", "Root")
	i := w.Root()

	assert.Equal(t, "Root", i.Title())

	assertChildrenListEmpty(t, i)
	assertItemDetached(t, i)
}

func TestItemDetach(t *testing.T) {
	t.Run("DetachFromMiddle", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		parent.Append(c)

		b.Detach()

		assertChildrenOrder(t, parent, a, c)
		assertItemDetached(t, b)
	})

	t.Run("DetachLast", func(t *testing.T) {
		parent, a, _, _ := newTestItems()
		parent.Append(a)

		a.Detach()

		assertChildrenListEmpty(t, parent)
		assertItemDetached(t, a)
	})
}

func TestItemMoveAbove(t *testing.T) {
	t.Run("MoveAboveHead", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		c.MoveAbove(a)

		assertChildrenOrder(t, parent, c, a, b)
	})

	t.Run("MoveAboveMiddle", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		c.MoveAbove(b)

		assertChildrenOrder(t, parent, a, c, b)
	})
}

func TestItemMoveBelow(t *testing.T) {
	t.Run("MoveBelowTail", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		c.MoveBelow(b)

		assertChildrenOrder(t, parent, a, b, c)
	})

	t.Run("MoveBelowMiddle", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		c.MoveBelow(a)

		assertChildrenOrder(t, parent, a, c, b)
	})
}

func TestItemMoveUp(t *testing.T) {
	t.Run("NilPrev", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		parent.Append(c)

		a.MoveUp()

		assertChildrenOrder(t, parent, a, b, c)
	})

	t.Run("NonNilPrev", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		parent.Append(c)

		b.MoveUp()

		assertChildrenOrder(t, parent, b, a, c)
	})
}

func TestItemMoveDown(t *testing.T) {
	t.Run("NilNext", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		parent.Append(c)

		c.MoveDown()

		assertChildrenOrder(t, parent, a, b, c)
	})

	t.Run("NonNilNext", func(t *testing.T) {
		parent, a, b, c := newTestItems()

		parent.Append(a)
		parent.Append(b)
		parent.Append(c)

		b.MoveDown()

		assertChildrenOrder(t, parent, a, c, b)
	})
}

func TestItemDemote(t *testing.T) {
	t.Run("NilPrev", func(t *testing.T) {
		parent, a, _, _ := newTestItems()
		parent.Append(a)

		a.Demote()

		assertChildrenOrder(t, parent, a)
	})

	t.Run("NonNilPrev", func(t *testing.T) {
		parent, a, b, c := newTestItems()
		parent.Append(a)
		parent.Append(b)
		parent.Append(c)

		a.SetCollapsed(true, false)
		b.Demote()

		assertChildrenOrder(t, parent, a, c)
		assertChildrenOrder(t, a, b)

		assert.False(t, a.Collapsed())
	})
}

func TestItemPromote(t *testing.T) {
	t.Run("RootItem", func(t *testing.T) {
		parent, _, _, _ := newTestItems()

		parent.Promote()

		assert.Nil(t, parent.Parent())
	})

	t.Run("SubRootItem", func(t *testing.T) {
		parent, a, _, _ := newTestItems()
		parent.Append(a)

		a.Promote()

		assertChildrenOrder(t, parent, a)
	})

	t.Run("NonNilParent", func(t *testing.T) {
		w := data.NewWorkspace("", "Grandparent")
		grandparentItem := w.Root()

		parentItem := w.NewItem("Parent")
		grandparentItem.Append(parentItem)

		childItem := w.NewItem("Child")
		parentItem.Append(childItem)

		childItem.Promote()

		assertChildrenListEmpty(t, childItem)
		assertChildrenListEmpty(t, parentItem)
		assertChildrenOrder(t, grandparentItem, parentItem, childItem)
	})
}

func TestItemPrepend(t *testing.T) {
	t.Run("EmptyList", func(t *testing.T) {
		parent, a, _, _ := newTestItems()

		parent.Prepend(a)

		assertChildrenOrder(t, parent, a)
	})

	t.Run("NonEmptyList", func(t *testing.T) {
		parent, a, b, _ := newTestItems()

		parent.Prepend(a)
		parent.Prepend(b)

		assertChildrenOrder(t, parent, b, a)
	})
}

func TestItemAppend(t *testing.T) {
	t.Run("EmptyList", func(t *testing.T) {
		parent, a, _, _ := newTestItems()

		parent.Append(a)

		assertChildrenOrder(t, parent, a)
	})

	t.Run("NonEmptyList", func(t *testing.T) {
		parent, a, b, _ := newTestItems()

		parent.Append(a)
		parent.Append(b)

		assertChildrenOrder(t, parent, a, b)
	})
}

func assertChildrenOrder(t *testing.T, parent *data.Item, children ...*data.Item) {
	t.Helper()

	for idx, child := range children {
		assert.Same(t, parent, child.Parent(), "item %s does not belong to parent %s", child.Title(), parent.Title())

		switch idx {
		case 0:
			assert.Same(t, child, parent.Head(), "item %s is not the children list head", child.Title())
			assert.Nil(t, child.Prev(), "item %s has non-nil previous sibling", child.Title())
		case len(children) - 1:
			assert.Same(t, child, parent.Tail(), "item %s in not the children list tail", child.Title())
			assert.Nil(t, child.Next(), "item %s has a non-nil next sibling", child.Title())
		default:
			assert.Same(t, children[idx-1], child.Prev(), "item %s has wrong previous sibling", child.Title())
			assert.Same(t, children[idx+1], child.Next(), "item %s has wrong next sibling", child.Title())
		}
	}
}

func assertChildrenListEmpty(t *testing.T, i *data.Item) {
	t.Helper()

	assert.Nil(t, i.Head())
	assert.Nil(t, i.Tail())
}

func newTestItems() (*data.Item, *data.Item, *data.Item, *data.Item) {
	w := data.NewWorkspace("", "Parent")
	a := w.NewItem("ChildA")
	b := w.NewItem("ChildB")
	c := w.NewItem("ChildC")

	return w.Root(), a, b, c
}

func assertItemDetached(t *testing.T, item *data.Item) {
	t.Helper()

	assert.Nil(t, item.Parent())
	assert.Nil(t, item.Prev())
	assert.Nil(t, item.Next())
}
