package main

type mselect struct {
	beginX       int
	beginY       int
	endX         int
	endY         int
	hasSelection bool
	isSelecting  bool
	isColumnMode bool
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
	}
	return &m
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
	// println("end...", x, y)
	m.isSelecting = false
	m.endX = x
	m.endY = y
}

func (m *mselect) clear() {
	m.hasSelection = false
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
