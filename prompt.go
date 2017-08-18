package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nsf/termbox-go"
)

type completion struct {
	text     string
	display  string
	startPos int
}

type box struct {
	completions []*completion
	selected    int
}

type prompt struct {
	prefix    string
	text      []rune
	pos       int
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
		fg:        fg,
		bg:        bg,
		completer: nil,
		popup:     box{nil, 1},
	}
	return &p
}

const cPopupFg = 234
const cPopupBg = 252
const cPopupSelectFg = 255
const cPopupSelectBg = 246

func (b *box) draw(x, y, w, h int) {
	var fg, bg termbox.Attribute

	n := len(b.completions)
	if n == 0 || y-n < 2 {
		return
	}
	width := b.getWidth()
	xoffset := x + width + 2 - w
	if xoffset > 0 {
		x -= xoffset
	}
	y -= n

	for j, compl := range b.completions {
		runes := []rune(fmt.Sprintf(" %-*s ", width, compl.display))
		if j == b.selected {
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

func (p *prompt) draw(x, y, w, h int, fg, bg termbox.Attribute) {
	for _, c := range p.prefix {
		termbox.SetCell(x, y, c, p.fg, p.bg)
		x++
	}
	/* x is now updated including the prefix */

	for i, r := range p.text {
		termbox.SetCell(x+i, y, r, fg, bg)
	}
	termbox.SetCursor(x+p.pos, y)
	p.popup.draw(x+p.pos, y, w, h)
}

func (p *prompt) insertRune(r rune) {
	p.text = append(p.text, '0')
	copy(p.text[p.pos+1:], p.text[p.pos:])
	p.text[p.pos] = r
	p.pos++
	p.getCompletions()
}

func (p *prompt) getCompletions() {
	if outPrompt.completer != nil {
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
	case termbox.KeyDelete:
		p.deleteRuneAtCursor()
	case termbox.KeyEsc:
		p.deleteAllRunes()
	case termbox.KeyArrowLeft:
		p.moveCursorBackward()
	case termbox.KeyArrowRight:
		p.moveCursorForward()
	case termbox.KeyHome:
		p.moveCursorToBegin()
	case termbox.KeyEnd:
		p.moveCursorToEnd()
	case termbox.KeySpace:
		p.insertRune(' ')
	default:
		if ev.Ch != 0 {
			p.insertRune(ev.Ch)
		}
	}
}
