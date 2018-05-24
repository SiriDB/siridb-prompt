package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nsf/termbox-go"
)

const cPopupFg = termbox.ColorBlack
const cPopupBg = termbox.ColorWhite
const cPopupSelectFg = termbox.ColorWhite
const cPopupSelectBg = termbox.ColorMagenta

type completion struct {
	text     string
	display  string
	startPos int
}

type box struct {
	completions []*completion
	selected    int
	iStart      int
	hidden      bool
}

type prompt struct {
	prefix    string
	text      []rune
	pos       int
	offset    int
	hideText  bool
	fg        termbox.Attribute
	bg        termbox.Attribute
	completer func(p *prompt) []*completion
	popup     box
}

func newPrompt(prefix string, fg, bg termbox.Attribute) *prompt {
	p := prompt{
		prefix:    prefix,
		text:      make([]rune, 0),
		pos:       0,
		offset:    0,
		fg:        fg,
		bg:        bg,
		hideText:  false,
		completer: nil,
		popup: box{
			completions: nil,
			selected:    -1,
			hidden:      false,
		},
	}
	return &p
}

func (b *box) draw(x, y, w, h int) {
	var fg, bg termbox.Attribute

	n := len(b.completions)
	if n == 0 || b.hidden == true {
		return
	}

	width := b.getWidth()
	xoffset := x + width + 2 - w
	if xoffset > 0 {
		x -= xoffset
	}

	start := 0
	end := n
	m := min(n, y-1)
	if n > m {
		if b.selected >= m {
			start = b.selected - m + 1
		}
		end = m + start
	}
	y -= m

	for j, compl := range b.completions[start:end] {
		runes := []rune(fmt.Sprintf(" %-*s ", width, compl.display))
		if j+start == b.selected {
			fg = cPopupSelectFg
			bg = cPopupSelectBg
		} else {
			fg = cPopupFg
			bg = cPopupBg
		}
		for i, r := range runes {
			termbox.SetCell(x+i, y+j, r, fg, bg)
		}
	}
}

func (b *box) getWidth() int {
	width := 0
	for _, compl := range b.completions {
		width = max(width, len(compl.display))
	}
	return width
}

func (b *box) selectNext() {
	n := len(b.completions)
	if n > 0 {
		b.selected++
		if b.selected >= n {
			b.selected = 0
		}
	}
}

func (b *box) selectPrev() {
	n := len(b.completions)
	if n > 0 {
		b.selected--
		if b.selected < 0 {
			b.selected = n - 1
		}
	}
}

func (p *prompt) insertSelected() {
	if p.popup.selected >= 0 {
		c := p.popup.completions[p.popup.selected]
		runes := []rune(c.text[c.startPos:])
		p.text = append(p.text[:p.popup.iStart], append(runes, p.text[p.pos:]...)...)
		p.pos = p.popup.iStart + len(runes)
	}
}

func (p *prompt) setText(s string) {
	p.showPopup()
	p.text = []rune(s)
	p.pos = len(p.text)
}

func (p *prompt) hidePopup() {
	p.popup.hidden = true
}

func (p *prompt) showPopup() {
	p.popup.hidden = false
}

func (p *prompt) draw(x, y, w, h int, fg, bg termbox.Attribute) {
	for _, c := range p.prefix {
		termbox.SetCell(x, y, c, p.fg, p.bg)
		x++
	}

	n := len(p.text)
	if p.offset+p.pos+x >= w {
		p.offset = min(x+n-w+1, p.pos)
	} else if p.offset > p.pos {
		p.offset = p.pos
	}
	n = min(p.offset+(w-x), n)
	if p.offset < 0 {
		p.offset = 0
	}
	for i, r := range p.text[p.offset:n] {
		if p.hideText {
			r = '*'
		}
		termbox.SetCell(x+i, y, r, fg, bg)
	}
	termbox.SetCursor(x+p.pos-p.offset, y)
	p.popup.draw(x+p.popup.iStart, y, w, h)
}

func (p *prompt) insertRune(r rune) {
	p.showPopup()
	p.text = append(p.text, '0')
	copy(p.text[p.pos+1:], p.text[p.pos:])
	p.text[p.pos] = r
	p.pos++
}

func (p *prompt) getCompletions() {
	if outPrompt.completer != nil {
		p.popup.selected = -1
		p.popup.iStart = p.pos
		p.popup.completions = outPrompt.completer(p)
		sort.Slice(p.popup.completions, func(i, j int) bool {
			return strings.Compare(p.popup.completions[i].display, p.popup.completions[j].display) < 0
		})
	}
}

func (p *prompt) clearCompletions() {
	for i := 0; i < len(p.popup.completions); i++ {
		p.popup.completions[i] = nil
	}
	p.popup.completions = p.popup.completions[:0]
}

func (p *prompt) hasCompletions() bool {
	return len(p.popup.completions) > 0
}

func (p *prompt) deleteRuneBeforeCursor() {
	if p.pos > 0 {
		p.text = append(p.text[:p.pos-1], p.text[p.pos:]...)
		p.pos--
	}
}

func (p *prompt) deleteRuneAtCursor() {
	if p.pos < len(p.text) {
		p.text = append(p.text[:p.pos], p.text[p.pos+1:]...)
	}
}

func (p *prompt) deleteAllRunes() {
	p.text = p.text[:0]
	p.pos = 0
}

func (p *prompt) moveCursorBackward() {
	if p.pos > 0 {
		p.pos--
	}
}

func (p *prompt) moveCursorForward() {
	if p.pos < len(p.text) {
		p.pos++
	}
}

func (p *prompt) moveCursorToEnd() {
	p.pos = len(p.text)
}

func (p *prompt) moveCursorToBegin() {
	p.pos = 0
}

func (p *prompt) textBeforeCursor() string {
	return string(p.text[:p.pos])
}

func (p *prompt) parse(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		p.deleteRuneBeforeCursor()
		p.getCompletions()
	case termbox.KeyDelete:
		p.deleteRuneAtCursor()
		p.clearCompletions()
	case termbox.KeyEsc:
		p.deleteAllRunes()
		p.clearCompletions()
	case termbox.KeyArrowLeft:
		p.moveCursorBackward()
		p.clearCompletions()
	case termbox.KeyArrowRight:
		p.moveCursorForward()
		p.clearCompletions()
	case termbox.KeyHome:
		p.moveCursorToBegin()
		p.clearCompletions()
	case termbox.KeyEnd:
		p.moveCursorToEnd()
		p.clearCompletions()
	case termbox.KeySpace:
		p.insertRune(' ')
		p.getCompletions()
	case termbox.KeyTab:
		if !p.hasCompletions() {
			p.getCompletions()
		}
		p.popup.selectNext()
		p.insertSelected()
	case termbox.KeyArrowDown:
		p.popup.selectNext()
		p.insertSelected()
	case termbox.KeyArrowUp:
		p.popup.selectPrev()
		p.insertSelected()
	default:
		if ev.Ch != 0 {
			p.insertRune(ev.Ch)
			p.getCompletions()
		}
	}
}
