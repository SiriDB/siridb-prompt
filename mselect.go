package main

import (
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

type mselect struct {
	beginX       int
	beginY       int
	endX         int
	endY         int
	hasSelection bool
	isSelecting  bool
	isColumnMode bool
	selection    []rune
}

func newMselect() *mselect {
	m := mselect{
		beginX:       0,
		beginY:       0,
		endX:         0,
		endY:         0,
		hasSelection: false,
		isSelecting:  false,
		isColumnMode: false,
		selection:    make([]rune, 0),
	}
	return &m
}

func (m *mselect) selectWord(x, y int) {
	w, _ := termbox.Size()
	offset := y * w
	cells := termbox.CellBuffer()
	for i := x; i >= 0; i-- {
		c := cells[offset+i]
		if unicode.IsSpace(c.Ch) {
			if i == x {
				return
			}
			break
		}
		m.hasSelection = true
		m.beginX = i
	}

	for m.endX = x + 1; m.endX < w; m.endX++ {
		c := cells[offset+m.endX]
		if unicode.IsSpace(c.Ch) {
			break
		}
	}
}

func (m *mselect) getSelection() []rune {
	return m.selection
}

func (m *mselect) setSelection() {
	if !m.hasSelection {
		return
	}
	w, _ := termbox.Size()
	lastY := -1
	m.selection = m.selection[:0]
	for i, cell := range termbox.CellBuffer() {
		x := i % w
		y := i / w
		if m.isInSelection(x, y) {
			if lastY != y {
				if lastY != -1 {
					m.selection = append(m.selection, '\n')
				}
				lastY = y
			}
			m.selection = append(m.selection, cell.Ch)
		}
	}
}

func (m *mselect) start(x int, y int) {
	m.endX = x
	m.endY = y
	if m.isSelecting {
		return
	}
	m.isSelecting = true
	m.beginX = x
	m.beginY = y
}

func (m *mselect) end(x int, y int) {
	m.isSelecting = false
	if m.beginX == x && m.beginY == y {
		m.selectWord(x, y)
	} else {
		m.endX = x
		m.endY = y
		m.hasSelection = true
	}
	m.setSelection()
}

func (m *mselect) clear() {
	m.hasSelection = false
	m.selection = m.selection[:0]
}

func isInSelection(x, y, bx, by, ex, ey int, isColumnMode bool) bool {
	if isColumnMode {
		return x >= bx && x <= ex && y >= by && y <= ey
	}

	if y < by || y > ey {
		return false
	}

	if (y > by && y < ey) || ((x >= bx || y != by) && (x <= ex || y != ey)) {
		return true
	}

	return false
}

func (m *mselect) isInSelectionColumnMode(x int, y int) bool {
	bx := m.beginX
	ex := m.endX
	by := m.beginY
	ey := m.endY
	if bx > ex {
		bx, ex = ex, bx
	}
	if by > ey {
		by, ey = ey, by
	}

	return x >= bx && x < ex && y >= by && y <= ey
}

func (m *mselect) isInSelectionNormalMode(x int, y int) bool {
	bx := m.beginX
	ex := m.endX
	by := m.beginY
	ey := m.endY

	if by > ey || (by == ey && bx > ex) {
		bx, ex = ex, bx
		by, ey = ey, by
	}

	if y < by || y > ey {
		return false
	}

	if (y > by && y < ey) || ((x >= bx || y != by) && (x < ex || y != ey)) {
		return true
	}

	return false
}

func (m *mselect) isInSelection(x int, y int) bool {
	if m.isColumnMode {
		return m.isInSelectionColumnMode(x, y)
	}
	return m.isInSelectionNormalMode(x, y)
}

func (m *mselect) draw(w, h int) {
	if !m.hasSelection && !m.isSelecting {
		return
	}
	for i, cell := range termbox.CellBuffer() {
		x := i % w
		y := i / w
		if m.isInSelection(x, y) {
			termbox.CellBuffer()[i] = termbox.Cell{
				Ch: cell.Ch,
				Fg: cell.Fg,
				Bg: colsel,
			}
		}
	}
}
