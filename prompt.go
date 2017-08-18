package main

import (
	"github.com/nsf/termbox-go"
)

type completion struct {
	text     string
	display  string
	startPos int
}

type box struct {
	completions []completion
	idx         int
}

type prompt struct {
	prefix    string
	text      []rune
	pos       int
	fg        termbox.Attribute
	bg        termbox.Attribute
	completer func(p *prompt) []completion
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
		popup:     box{nil, -1},
	}
	return &p
}

func (p *prompt) draw(x, y, w int, fg, bg termbox.Attribute) {
	for _, c := range p.prefix {
		termbox.SetCell(x, y, c, p.fg, p.bg)
		x++
	}
	/* x is now updated including the prefix */

	for i, r := range p.text {
		termbox.SetCell(x+i, y, r, fg, bg)
	}

	termbox.SetCursor(x+p.pos, y)
}

func (p *prompt) insertRune(r rune) {
	p.text = append(p.text, '0')
	copy(p.text[p.pos+1:], p.text[p.pos:])
	p.text[p.pos] = r
	p.pos++
	if outPrompt.completer != nil {
		p.popup.completions = outPrompt.completer(p)
		if len(p.popup.completions) {

		}
	}
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
