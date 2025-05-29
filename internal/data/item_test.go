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
	"github.com/stretchr/testify/require"

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
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		root.Append(c)

		b.Detach()

		assertChildrenOrder(t, root, a, c)
		assertItemDetached(t, b)
	})

	t.Run("DetachLast", func(t *testing.T) {
		w, a, _, _ := newTestItems()
		root := w.Root()

		root.Append(a)

		a.Detach()

		assertChildrenListEmpty(t, root)
		assertItemDetached(t, a)
	})
}

func TestItemMoveAbove(t *testing.T) {
	t.Run("MoveAboveHead", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		c.MoveAbove(a)

		assertChildrenOrder(t, root, c, a, b)
	})

	t.Run("MoveAboveMiddle", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		c.MoveAbove(b)

		assertChildrenOrder(t, root, a, c, b)
	})
}

func TestItemMoveBelow(t *testing.T) {
	t.Run("MoveBelowTail", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		c.MoveBelow(b)

		assertChildrenOrder(t, root, a, b, c)
	})

	t.Run("MoveBelowMiddle", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		c.MoveBelow(a)

		assertChildrenOrder(t, root, a, c, b)
	})
}

func TestItemMoveUp(t *testing.T) {
	t.Run("NilPrev", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		root.Append(c)

		a.MoveUp()

		assertChildrenOrder(t, root, a, b, c)
	})

	t.Run("NonNilPrev", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		root.Append(c)

		b.MoveUp()

		assertChildrenOrder(t, root, b, a, c)
	})
}

func TestItemMoveDown(t *testing.T) {
	t.Run("NilNext", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		root.Append(c)

		c.MoveDown()

		assertChildrenOrder(t, root, a, b, c)
	})

	t.Run("NonNilNext", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		root.Append(c)

		b.MoveDown()

		assertChildrenOrder(t, root, a, c, b)
	})
}

func TestItemDemote(t *testing.T) {
	t.Run("NilPrev", func(t *testing.T) {
		w, a, _, _ := newTestItems()
		root := w.Root()

		root.Append(a)

		a.Demote()

		assertChildrenOrder(t, root, a)
	})

	t.Run("NonNilPrev", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)
		root.Append(c)

		a.SetCollapsed(true, false)
		b.Demote()

		assertChildrenOrder(t, root, a, c)
		assertChildrenOrder(t, a, b)

		assert.False(t, a.Collapsed())
	})
}

func TestItemPromote(t *testing.T) {
	t.Run("RootItem", func(t *testing.T) {
		w, _, _, _ := newTestItems()
		root := w.Root()

		root.Promote()

		assert.Nil(t, root.Parent())
	})

	t.Run("SubRootItem", func(t *testing.T) {
		w, a, _, _ := newTestItems()
		root := w.Root()

		root.Append(a)

		a.Promote()

		assertChildrenOrder(t, root, a)
	})

	t.Run("NonNilParent", func(t *testing.T) {
		w, a, b, _ := newTestItems()
		root := w.Root()

		root.Append(a)
		a.Append(b)

		b.Promote()

		assertChildrenListEmpty(t, b)
		assertChildrenListEmpty(t, a)
		assertChildrenOrder(t, root, a, b)
	})
}

func TestItemPrepend(t *testing.T) {
	t.Run("EmptyList", func(t *testing.T) {
		w, a, _, _ := newTestItems()
		root := w.Root()

		root.Prepend(a)

		assertChildrenOrder(t, root, a)
	})

	t.Run("NonEmptyList", func(t *testing.T) {
		w, a, b, _ := newTestItems()
		root := w.Root()

		root.Prepend(a)
		root.Prepend(b)

		assertChildrenOrder(t, root, b, a)
	})
}

func TestItemAppend(t *testing.T) {
	t.Run("EmptyList", func(t *testing.T) {
		w, a, _, _ := newTestItems()
		root := w.Root()

		root.Append(a)

		assertChildrenOrder(t, root, a)
	})

	t.Run("NonEmptyList", func(t *testing.T) {
		w, a, b, _ := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)

		assertChildrenOrder(t, root, a, b)
	})
}

func TestItemDepth(t *testing.T) {
	w, a, b, c := newTestItems()
	root := w.Root()

	root.Append(a)
	a.Append(b)
	b.Append(c)

	assert.Equal(t, 0, root.Depth())
	assert.Equal(t, 1, a.Depth())
	assert.Equal(t, 2, b.Depth())
	assert.Equal(t, 3, c.Depth())

	w.SetRoot(a)

	assert.Equal(t, -1, root.Depth())
	assert.Equal(t, 0, a.Depth())
	assert.Equal(t, 1, b.Depth())
	assert.Equal(t, 2, c.Depth())

	w.SetRoot(b)

	assert.Equal(t, -1, root.Depth())
	assert.Equal(t, -1, a.Depth())
	assert.Equal(t, 0, b.Depth())
	assert.Equal(t, 1, c.Depth())
}

func TestItemRealRoot(t *testing.T) {
	w, a, b, c := newTestItems()
	root := w.Root()

	root.Append(a)
	a.Append(b)
	b.Append(c)

	assert.Same(t, root, root.RealRoot())
	assert.Same(t, root, a.RealRoot())
	assert.Same(t, root, b.RealRoot())
	assert.Same(t, root, c.RealRoot())

	w.SetRoot(a)

	assert.Same(t, root, root.RealRoot())
	assert.Same(t, root, a.RealRoot())
	assert.Same(t, root, b.RealRoot())
	assert.Same(t, root, c.RealRoot())

	w.SetRoot(b)

	assert.Same(t, root, root.RealRoot())
	assert.Same(t, root, a.RealRoot())
	assert.Same(t, root, b.RealRoot())
	assert.Same(t, root, c.RealRoot())
}

func TestItemSetCollapsed(t *testing.T) {
	t.Run("SetToTrue", func(t *testing.T) {
		t.Run("No children", func(t *testing.T) {
			w, _, _, _ := newTestItems()
			root := w.Root()

			assert.False(t, root.Collapsed())

			root.SetCollapsed(true, false)

			assert.False(t, root.Collapsed())
		})

		t.Run("Non-recursive", func(t *testing.T) {
			w, a, b, c := newTestItems()
			root := w.Root()

			root.Append(a)
			a.Append(b)
			b.Append(c)

			root.SetCollapsed(true, false)

			assert.True(t, root.Collapsed())
			assert.False(t, a.Collapsed())
			assert.False(t, b.Collapsed())
			assert.False(t, c.Collapsed())
		})

		t.Run("Recursive", func(t *testing.T) {
			w, a, b, c := newTestItems()
			root := w.Root()

			root.Append(a)
			a.Append(b)
			b.Append(c)

			root.SetCollapsed(true, true)

			assert.True(t, root.Collapsed())
			assert.True(t, a.Collapsed())
			assert.True(t, b.Collapsed())
			assert.False(t, c.Collapsed())
		})
	})
	t.Run("SetToFalse", func(t *testing.T) {
		t.Run("No children", func(t *testing.T) {
			w, _, _, _ := newTestItems()
			root := w.Root()

			assert.False(t, root.Collapsed())

			root.SetCollapsed(false, false)

			assert.False(t, root.Collapsed())
		})

		t.Run("Non-recursive", func(t *testing.T) {
			w, a, b, c := newTestItems()
			root := w.Root()

			root.Append(a)
			a.Append(b)
			b.Append(c)

			// prepare tested state
			root.SetCollapsed(true, true)

			root.SetCollapsed(false, false)

			assert.False(t, root.Collapsed())
			assert.True(t, a.Collapsed())
			assert.True(t, b.Collapsed())
			assert.False(t, c.Collapsed())
		})

		t.Run("Recursive", func(t *testing.T) {
			w, a, b, c := newTestItems()
			root := w.Root()

			root.Append(a)
			a.Append(b)
			b.Append(c)

			// prepare tested state
			root.SetCollapsed(true, true)

			root.SetCollapsed(false, true)

			assert.False(t, root.Collapsed())
			assert.False(t, a.Collapsed())
			assert.False(t, b.Collapsed())
			assert.False(t, c.Collapsed())
		})
	})
}

func TestItemDisplayChildren(t *testing.T) {
	t.Run("EmptyParent", func(t *testing.T) {
		w, _, _, _ := newTestItems()
		root := w.Root()

		children := root.DisplayedChildren()
		require.Empty(t, children)
	})

	t.Run("SimpleCase", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		a.Append(b)
		b.Append(c)

		children := root.DisplayedChildren()
		require.Len(t, children, 3)
		assert.Same(t, a, children[0])
		assert.Same(t, b, children[1])
		assert.Same(t, c, children[2])
	})

	t.Run("CollapsedChild", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		a.Append(b)
		b.Append(c)

		b.SetCollapsed(true, false)

		children := root.DisplayedChildren()
		require.Len(t, children, 2)
		assert.Same(t, a, children[0])
		assert.Same(t, b, children[1])
	})
}

func TestItemPrevRow(t *testing.T) {
	t.Run("No previous sibling", func(t *testing.T) {
		t.Run("Parent is root", func(t *testing.T) {
			w, a, _, _ := newTestItems()
			root := w.Root()

			root.Append(a)

			prevRow := a.PrevRow()
			assert.Nil(t, prevRow)
		})

		t.Run("Parent is not root", func(t *testing.T) {
			w, a, b, _ := newTestItems()

			root := w.Root()

			root.Append(a)
			a.Append(b)

			prevRow := b.PrevRow()
			assert.Same(t, a, prevRow)
		})
	})

	t.Run("With previous sibling", func(t *testing.T) {
		t.Run("Previous sibling has no children", func(t *testing.T) {
			w, a, b, _ := newTestItems()
			root := w.Root()

			root.Append(a)
			root.Append(b)

			prevRow := b.PrevRow()
			assert.Same(t, a, prevRow)
		})

		t.Run("Previour sibling has children", func(t *testing.T) {
			t.Run("Previous sibling is collapsed", func(t *testing.T) {
				w, a, b, c := newTestItems()
				root := w.Root()

				root.Append(a)
				a.Append(b)

				a.SetCollapsed(true, false)

				root.Append(c)

				prevRow := c.PrevRow()
				assert.Same(t, a, prevRow)
			})

			t.Run("Previous sibling is not collapsed", func(t *testing.T) {
				w, a, b, c := newTestItems()
				root := w.Root()

				root.Append(a)
				a.Append(b)

				root.Append(c)

				prevRow := c.PrevRow()
				assert.Same(t, b, prevRow)
			})
		})
	})
}

func TestItemNextRow(t *testing.T) {
	t.Run("Item has children and is not collapsed", func(t *testing.T) {
		w, a, b, _ := newTestItems()
		root := w.Root()

		root.Append(a)
		a.Append(b)

		nextRow := a.NextRow()
		assert.Same(t, b, nextRow)
	})

	t.Run("Item has next sibling", func(t *testing.T) {
		w, a, b, _ := newTestItems()
		root := w.Root()

		root.Append(a)
		root.Append(b)

		nextRow := a.NextRow()
		assert.Same(t, b, nextRow)
	})

	t.Run("Item's parent has next sibling", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		a.Append(b)
		root.Append(c)

		nextRow := b.NextRow()
		assert.Same(t, c, nextRow)
	})

	t.Run("No next row", func(t *testing.T) {
		w, a, _, _ := newTestItems()
		root := w.Root()

		root.Append(a)

		nextRow := a.NextRow()
		assert.Nil(t, nextRow)
	})

	t.Run("No next row while zoomed", func(t *testing.T) {
		w, a, b, c := newTestItems()
		root := w.Root()

		root.Append(a)
		a.Append(b)
		root.Append(c)

		w.SetRoot(a)

		nextRow := b.NextRow()
		assert.Nil(t, nextRow)
	})
}

func newTestItems() (*data.Workspace, *data.Item, *data.Item, *data.Item) {
	w := data.NewWorkspace("", "Parent")

	a := w.NewItem("ChildA")
	b := w.NewItem("ChildB")
	c := w.NewItem("ChildC")

	return w, a, b, c
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

func assertItemDetached(t *testing.T, item *data.Item) {
	t.Helper()

	assert.Nil(t, item.Parent())
	assert.Nil(t, item.Prev())
	assert.Nil(t, item.Next())
}
