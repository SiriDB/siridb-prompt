package main

import (
	"fmt"
	"sort"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

const cModePretty = 0
const cModeJSON = 1

type view struct {
	lines []string
	query *query
	mode  int
	pos   int
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func newView() *view {
	v := view{
		lines: make([]string, 0),
		query: nil,
		mode:  cModePretty,
		pos:   -1,
	}
	return &v
}

func (v *view) setModeJSON() {
	v.mode = cModeJSON
}

func (v *view) addString(s string, w int) {
	r := []rune(s)

	for i := 0; i < len(r); i += w {
		batch := r[i:min(i+w, len(r))]
		v.lines = append(v.lines, string(batch))
	}
}

func (v *view) newLine() {
	v.lines = append(v.lines, "")
}

func (v *view) addTable(columns []string, rows []interface{}, limit, mustSort bool, w int) error {
	nc := len(columns)
	if nc == 0 {
		return fmt.Errorf("cannot make table: : at least one column is required")
	}
	nr := len(rows)
	if limit {
		nr = min(nr, 100)
	}
	colSizes := make([]int, nc)
	colFormatters := make([]func(i interface{}) string, nc)
	fmtCols := make([]string, nc)
	fmtRows := make([][]string, nr)
	for i, s := range columns {
		colSizes[i] = len(s)
		switch s {
		case "start", "end", "timestamp":
			colFormatters[i] = fmtTimestamp
		case "size", "buffer_size":
			colFormatters[i] = fmtSizeBytes
		case "received_points", "selected_points":
			colFormatters[i] = fmtLargeNum
		case "mem_usage":
			colFormatters[i] = fmtSizeMb
		case "uptime":
			colFormatters[i] = fmtDuration
		default:
			colFormatters[i] = func(i interface{}) string { return fmt.Sprint(i) }
		}
	}
	for r, rowIf := range rows {
		if row, ok := rowIf.([]interface{}); ok {
			if nc != len(row) {
				return fmt.Errorf("cannot make table: invalid row length found")
			}
			fmtRows[r] = make([]string, nc)
			for c, col := range row {
				s := colFormatters[c](col)
				if len(s) > colSizes[c] {
					colSizes[c] = len(s)
				}
				fmtRows[r][c] = s
			}
		}
		if r > nr {
			break
		}
	}
	if mustSort {
		sort.Slice(fmtRows, func(i, j int) bool {
			return strings.Compare(fmtRows[i][0], fmtRows[j][0]) < 0
		})
	}
	for i, w := range colSizes {
		fmtCols[i] = fmt.Sprintf("%-*s", w, columns[i])
	}
	v.addString(strings.Join(fmtCols, " \u2503 "), w)
	for i, w := range colSizes {
		fmtCols[i] = strings.Repeat("\u2501", w)
	}
	v.addString(strings.Join(fmtCols, "\u2501\u2547\u2501"), w)
	for _, row := range fmtRows {
		for c, w := range colSizes {
			fmtCols[c] = fmt.Sprintf("%-*s", w, row[c])
		}
		v.addString(strings.Join(fmtCols, " \u2502 "), w)
	}
	return nil
}

func (v *view) addPretty(w int) error {
	m, ok := v.query.res.(map[string]interface{})
	if !ok {
		return fmt.Errorf("got an unexpected map")
	}

	if stop, err := v.tryList(m, w); stop {
		return err
	}
	return nil
}

func (v *view) tryList(m map[string]interface{}, w int) (bool, error) {
	var colIf interface{}
	var colArr []interface{}
	var columns []string

	var ok bool
	if colIf, ok = m["columns"]; !ok {
		return false, fmt.Errorf("columns not found")
	}
	if colArr, ok = colIf.([]interface{}); !ok {
		return false, fmt.Errorf("columns not a slice")
	}
	if len(colArr) == 0 {
		return false, fmt.Errorf("zero columns found")
	}

	columns = make([]string, len(colArr))

	for i, col := range colArr {
		if s, ok := col.(string); ok {
			columns[i] = s
		} else {
			return false, fmt.Errorf("columns contains non string")
		}
	}

	for k, data := range m {
		switch k {
		case "__timeit__", "columns":
			// skip
		default:
			var rows []interface{}
			var ok bool
			if rows, ok = data.([]interface{}); !ok {
				return true, fmt.Errorf("%s not a slice", k)
			}
			v.addTable(columns, rows, false, true, w)
		}
	}

	return true, nil
}

func (v *view) error(err error, w int) {
	logger.append(err.Error())
	v.addString(err.Error(), w)
}

func (v *view) append(q *query, w int) {
	v.query = q
	v.addString(fmt.Sprintf("> %s", q.req), w)
	if q.err != nil {
		v.error(q.err, w)
	} else if v.mode == cModeJSON {
		s, e := q.json()
		if e != nil {
			v.error(q.err, w)
			return
		}
		v.addString(s, w)
	} else if v.mode == cModePretty {
		err := v.addPretty(w)
		if err != nil {
			v.error(err, w)
		}
	}
	v.newLine()
	v.pos = -1
}

func (v *view) draw(w, h int) {
	y := 1
	n := len(v.lines)
	m := h - 2

	if n < m {
		v.pos = -1
	}

	if v.pos != -1 && v.pos < m {
		v.pos = m
	}

	pos := v.pos
	if pos == -1 {
		pos = n
	}
	start := pos - m
	if start < 0 {
		start = 0
	}

	for _, line := range v.lines[start:pos] {
		runes := []rune(line)
		for x, r := range runes {
			termbox.SetCell(x, y, r, coldef, coldef)
		}
		y++
	}
}

func (v *view) up() {
	if v.pos > 0 {
		v.pos--
	} else if v.pos == -1 {
		v.pos = len(v.lines) - 1
	}
}

func (v *view) down() {
	if v.pos != -1 {
		v.pos++
	}
	if v.pos >= len(v.lines) {
		v.pos = -1
	}
}

func (v *view) pageUp() {
	if v.pos == 0 {
		return
	}
	_, h := termbox.Size()

	if v.pos == -1 {
		v.pos = len(v.lines)
	}

	v.pos -= (h - 3)

	if v.pos < 0 {
		v.pos = 0
	}
}

func (v *view) pageDown() {
	if v.pos == -1 {
		return
	}
	_, h := termbox.Size()

	v.pos += (h - 3)

	if v.pos >= len(v.lines) {
		v.pos = -1
	}
}
