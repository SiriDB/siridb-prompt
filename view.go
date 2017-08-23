package main

import (
	"fmt"
	"sort"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

const cModePretty = 0
const cModeJSON = 1
const cLimitRows = 200

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
	if len(s) == 0 {
		v.newLine()
		return
	}

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
		return fmt.Errorf("cannot make table: at least one column is required")
	}
	nr := len(rows)
	if limit {
		nr = min(nr, cLimitRows)
	}
	colSizes := make([]int, nc)
	colFormatters := make([]func(i interface{}) string, nc)
	fmtCols := make([]string, nc)
	fmtRows := make([][]string, nr)
	for i, s := range columns {
		colSizes[i] = len(s)
		colFormatters[i] = getFormatter(s)
	}
	for r, rowIf := range rows {
		if r >= nr {
			break
		}
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
	if len(rows) > nr {
		v.addString("...", w)
	}
	return nil
}

func (v *view) addPretty(w int) error {
	m, ok := v.query.res.(map[string]interface{})
	if !ok {
		return fmt.Errorf("got an unexpected map")
	}

	if err := v.tryShow(m, w); err != nil {
		return err
	}

	if err := v.tryCalc(m, w); err != nil {
		return err
	}

	if err := v.tryMotd(m, w); err != nil {
		return err
	}

	if err := v.trySuccess(m, w); err != nil {
		return err
	}

	if err := v.tryError(m, w); err != nil {
		return err
	}

	if err := v.tryHelp(m, w); err != nil {
		return err
	}

	if err := v.tryCount(m, w); err != nil {
		return err
	}

	if columns := getListColumns(m); columns != nil {
		if err := v.tryList(m, w, columns); err != nil {
			return err
		}
	} else {
		if err := v.trySelect(m, w); err != nil {
			return err
		}
	}

	if err := v.tryTimeit(m, w); err != nil {
		return err
	}

	return nil
}

func (v *view) tryTimeit(m map[string]interface{}, w int) error {
	var err error
	if data, ok := m["__timeit__"]; ok {
		var columns = []string{"server", "time"}
		var rows []interface{}
		if items, ok := data.([]interface{}); ok {
			for _, item := range items {
				if itm, ok := item.(map[string]interface{}); ok {
					if server, ok := itm["server"]; ok {
						if t, ok := itm["time"]; ok {
							rows = append(rows, []interface{}{server, t})
						}
					}
				} else {
					break
				}
			}
		}
		if len(rows) > 0 {
			v.newLine()
			err = v.addTable(columns, rows, false, true, w)
		}
	}
	return err
}

func (v *view) tryShow(m map[string]interface{}, w int) error {
	var err error
	if data, ok := m["data"]; ok {
		var columns = []string{"name", "value"}
		var rows []interface{}
		if items, ok := data.([]interface{}); ok {
			for _, item := range items {
				if itm, ok := item.(map[string]interface{}); ok {
					if name, ok := itm["name"]; ok {
						if value, ok := itm["value"]; ok {
							n := fmt.Sprint(name)
							v := getFormatter(n)(value)
							rows = append(rows, []interface{}{n, v})
						}
					} else {
						break
					}
				} else {
					break
				}
			}
		}
		if len(rows) > 0 {
			err = v.addTable(columns, rows, false, true, w)
		}
	}
	return err
}

func (v *view) tryCalc(m map[string]interface{}, w int) error {
	if data, ok := m["calc"]; ok {
		n, ok := data.(int)
		if ok {
			v.addString(fmt.Sprint(n), w)
			v.addString(fmtTimestampUTC(n), w)
			v.addString(fmt.Sprintf("Local time: %s", fmtTimestamp(n)), w)
		}
	}
	return nil
}

func (v *view) tryMotd(m map[string]interface{}, w int) error {
	if data, ok := m["motd"]; ok {
		s, ok := data.(string)
		if ok {
			for _, line := range strings.Split(s, "\n") {
				v.addString(line, w)
			}
		}
	}
	return nil
}

func (v *view) trySuccess(m map[string]interface{}, w int) error {
	if data, ok := m["success_msg"]; ok {
		s, ok := data.(string)
		if ok {
			v.addString(s, w)
		}
	}
	return nil
}

func (v *view) tryError(m map[string]interface{}, w int) error {
	if data, ok := m["error_msg"]; ok {
		s, ok := data.(string)
		if ok {
			v.addString(s, w)
		}
	}
	return nil
}

func (v *view) tryHelp(m map[string]interface{}, w int) error {
	if data, ok := m["help"]; ok {
		s, ok := data.(string)
		if ok {
			for _, line := range strings.Split(s, "\n") {
				v.addString(line, w)
			}
		}
	}
	return nil
}

func (v *view) tryCount(m map[string]interface{}, w int) error {
	for k, data := range m {
		switch k {
		case "calc":
			// skip
		default:
			n, ok := data.(int)
			if ok {
				v.addString(fmtLargeNum(n), w)
			}
		}
	}
	return nil
}

func getListColumns(m map[string]interface{}) []string {
	if colIf, ok := m["columns"]; ok {
		var columns []string

		if colArr, ok := colIf.([]interface{}); ok {
			for _, col := range colArr {
				if s, ok := col.(string); ok {
					columns = append(columns, s)
				}
			}
		}

		if len(columns) > 0 {
			return columns
		}
	}
	return nil
}

func (v *view) tryList(m map[string]interface{}, w int, columns []string) error {
	var err error
	for k, data := range m {
		switch k {
		case "__timeit__", "columns":
			// skip
		default:
			var rows []interface{}
			var ok bool
			if rows, ok = data.([]interface{}); ok {
				err = v.addTable(columns, rows, false, true, w)
			}
		}
	}
	return err
}

func (v *view) trySelect(m map[string]interface{}, w int) error {
	var columns = []string{"timestamp", "value"}
	n := 0
	for series, data := range m {
		if points, ok := data.([]interface{}); ok {
			if len(points) > 0 {
				if point, test := points[0].([]interface{}); !test || len(point) != 2 {
					ok = false
				}
			}
			if ok {
				if n > cLimitRows {
					v.newLine()
					v.addString("...more series are found but hidden from output", w)
					break
				}
				n++
				v.newLine()
				v.addString(fmt.Sprintf("# %s", series), w)
				if err := v.addTable(columns, points, true, false, w); err != nil {
					return err
				}

			}
		}
	}
	return nil
}

func (v *view) error(err error, w int) {
	logger.ch <- err.Error()
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
