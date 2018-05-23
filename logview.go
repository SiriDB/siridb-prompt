package main

import (
	"fmt"
	"log"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type logEntry struct {
	dtime time.Time
	msg   string
}

type logView struct {
	pos        int
	entries    []logEntry
	ch         chan string
	useTermbox bool
	stderr     *log.Logger
}

func newLogView() *logView {
	l := logView{
		pos:        -1,
		entries:    make([]logEntry, 0),
		ch:         make(chan string),
		useTermbox: false,
		stderr:     nil,
	}
	return &l
}

func (l *logView) setMode(s string) {
	switch s {
	case "CONSOLE":
		l.useTermbox = false
	case "TERMBOX":
		l.useTermbox = true
	default:
		log.Fatal("Invalid mode: ", s)
	}
}

func (l *logView) handle() {
	for {
		msg := <-l.ch
		if l.stderr != nil {
			l.stderr.Println(msg)
		} else {
			l.entries = append(l.entries, logEntry{time.Now(), msg})
			if l.useTermbox {
				draw()
			}
		}
	}
}

func (l *logView) toStdErr() {
	l.stderr = log.New(os.Stderr, "", 0)
	for _, entry := range l.entries {
		l.stderr.Println(entry.msg)
	}
}

func (l *logView) draw(w, h int) {
	y := 1
	n := len(l.entries)
	m := h - 1

	if n < m {
		l.pos = -1
	}

	if l.pos != -1 && l.pos < m {
		l.pos = m
	}

	pos := l.pos
	if pos == -1 {
		pos = n
	}
	start := pos - m
	if start < 0 {
		start = 0
	}

	for _, entry := range l.entries[start:pos] {
		s := fmt.Sprintf("%s %s", entry.dtime.Format(time.UnixDate), entry.msg)
		for x, c := range s {
			termbox.SetCell(x, y, c, coldef, coldef)
		}
		y++

	}
	termbox.HideCursor()
}

func (l *logView) up() {
	if l.pos > 0 {
		l.pos--
	} else if l.pos == -1 {
		l.pos = len(l.entries) - 1
	}
}

func (l *logView) down() {
	if l.pos != -1 {
		l.pos++
	}
	if l.pos >= len(l.entries) {
		l.pos = -1
	}
}

func (l *logView) pageUp() {
	if l.pos == 0 {
		return
	}
	_, h := termbox.Size()

	if l.pos == -1 {
		l.pos = len(l.entries)
	}

	l.pos -= (h - 2)

	if l.pos < 0 {
		l.pos = 0
	}
}

func (l *logView) pageDown() {
	if l.pos == -1 {
		return
	}
	_, h := termbox.Size()

	l.pos += (h - 2)

	if l.pos >= len(l.entries) {
		l.pos = -1
	}
}
