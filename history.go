package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
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

	if err := os.MkdirAll(path.Dir(*h.fn), os.ModePerm); err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(*h.fn)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	f.WriteString(strings.Join(h.entries, "\n"))
}

func (h *history) append() {
	if h.fn == nil {
		return
	}

	if err := os.MkdirAll(path.Dir(*h.fn), os.ModePerm); err != nil {
		logger.ch <- err.Error()
	}

	f, err := os.OpenFile(*h.fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		logger.ch <- err.Error()
	} else {
		defer f.Close()
		s := his.entries[len(his.entries)-1]
		if len(his.entries) > 1 {
			f.WriteString(fmt.Sprintf("\n%s", s))
		} else {
			f.WriteString(s)
		}
	}
}

func (h *history) insert(s string) {
	if len(s) > 0 && (len(h.entries) == 0 || strings.Compare(s, h.entries[0]) != 0) {
		h.entries = append([]string{s}, h.entries...)
	}

	h.pos = -1

	if len(his.entries) > his.size {
		his.entries = his.entries[:his.size]
	}

	h.append()
}

func (h *history) prev() string {
	n := len(h.entries)
	if n == 0 {
		return ""
	}
	if h.pos < n-1 {
		h.pos++
	}
	return h.entries[h.pos]
}

func (h *history) next() string {
	if h.pos >= 0 {
		h.pos--
	}
	if h.pos == -1 {
		return ""
	}
	return h.entries[h.pos]
}
