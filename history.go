package main

import (
	"bufio"
	"os"
	"strings"
)

type history struct {
	entries []string
	pos     int
	size    int
	fn      *string
}

func newHistory(size int, fn *string) *history {
	h := history{
		entries: make([]string, 0),
		pos:     -1,
		fn:      fn,
		size:    size,
	}
	return &h
}

func (h *history) load() {
	if h.fn == nil {
		return
	}

	f, err := os.Open(*h.fn)
	if err != nil {
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for n := 0; n < h.size && scanner.Scan(); n++ {
		his.entries = append(his.entries, scanner.Text())
	}
}

func (h *history) save() {
	if h.fn == nil {
		return
	}

	f, err := os.Create(*h.fn)
	if err != nil {
		return
	}

	defer f.Close()

	f.WriteString(strings.Join(h.entries, "\n"))
}

func (h *history) insert(s string) {
	if len(s) > 0 && (len(h.entries) == 0 || strings.Compare(s, h.entries[len(h.entries)-1]) != 0) {
		h.entries = append([]string{s}, h.entries...)
		h.pos = len(h.entries)
	}

	if len(his.entries) > his.size {
		his.entries = his.entries[:his.size]
	}
}
