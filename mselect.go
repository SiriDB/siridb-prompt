package main

type mselect struct {
	beginX      int
	beginY      int
	endX        int
	endY        int
	hasSelection bool
	isSelecting bool
	isColumnMode bool
}

func newMselect() *mselect {
	m := mselect{
		beginX:      0,
		beginY:      0,
		endX:        0,
		endY:        0,
		hasSelection: false,
		isSelecting: false,
		isColumnMode: false,
	}
	return &m
}

func (m *mselect) clear() {
	hasSelection = false
}

func (m *mselect) isInSelectoin(x int, y, int) bool {
	if is
	if m.isColumnMode {
		return x >= m.beginX && x <= m.endX && y >= m.beginY && y <= m.endY
	}
	return x >= m.beginX && y >= m.beginY &&
}