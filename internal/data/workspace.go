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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	workspaceFilename = "workspace.xml"
	storageVersion    = 2

	xmlElemItem             = "item"
	xmlItemAttrStatus       = "status"
	xmlItemAttrCollapsed    = "collapsed"
	xmlItemAttrSelected     = "selected"
	xmlItemAttrZoomedIn     = "zoomed-in"
	xmlElemTitle            = "title"
	xmlElemWorkspace        = "oli-workspace"
	xmlWorkspaceAttrVersion = "version"
)

type Workspace struct {
	directory string

	root   *Item
	cursor *Item
}

func NewWorkspace(directory, rootTitle string) *Workspace {
	w := &Workspace{
		directory: directory,
	}
	w.root = w.NewItem(rootTitle)

	return w
}

func LoadWorkspace(directory string) (*Workspace, error) {
	p := filepath.Join(directory, workspaceFilename)
	w := NewWorkspace(directory, "Home")

	if _, err := os.Stat(p); os.IsNotExist(err) {
		i := w.NewItem("")
		w.root.Append(i)
		w.cursor = i

		return w, w.Save()
	} else if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return w, xml.Unmarshal(data, w)
}

// NewItem returns a new item not attached to any list.
func (w *Workspace) NewItem(title string) *Item {
	i := &Item{
		workspace: w,
		title:     title,
	}

	return i
}

func (w *Workspace) Root() *Item {
	return w.root
}

func (w *Workspace) SetRoot(item *Item) {
	w.root = item
}

func (w *Workspace) ZoomOut() {
	if w.root.parent == nil {
		return
	}

	if w.root.collapsed {
		w.cursor = w.root
	}
	w.root = w.root.parent
}

func (w *Workspace) Cursor() *Item {
	return w.cursor
}

func (w *Workspace) SetCursor(item *Item) {
	w.cursor = item
}

func (w *Workspace) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = xmlElemWorkspace
	start.Attr = []xml.Attr{
		{Name: xml.Name{Local: xmlWorkspaceAttrVersion}, Value: strconv.Itoa(storageVersion)},
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.Encode(w.root.RealRoot()); err != nil {
		return err
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (w *Workspace) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var vStr string
	for _, attr := range start.Attr {
		if attr.Name.Local == xmlWorkspaceAttrVersion {
			vStr = attr.Value
		}
	}

	if v, err := strconv.Atoi(vStr); err != nil {
		return fmt.Errorf("failed to parse storage version: %w", err)
	} else if v != storageVersion {
		return fmt.Errorf("unsupported storage version %d", v)
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
			case xmlElemItem:
				if err := d.DecodeElement(&w.root, &se); err != nil {
					return err
				}
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

func (w *Workspace) Save() error {
	p := filepath.Join(w.directory, workspaceFilename)
	if _, err := os.Stat(p); err == nil {
		backupFilename := fmt.Sprintf("%s.bak.%d", workspaceFilename, time.Now().Unix())
		backupPath := filepath.Join(w.directory, backupFilename)
		if err := os.Rename(p, backupPath); err != nil {
			return err
		}
	}

	data, err := xml.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, 0600)
}
