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
		mode:  cModeJSON,
		pos:   -1,
	}
	return &v
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
	fmtCols := make([]string, nc)
	fmtRows := make([][]string, nr)
	for i, s := range columns {
		colSizes[i] = len(s)
	}
	for r, rowIf := range rows {
		if row, ok := rowIf.([]interface{}); ok {
			if nc != len(row) {
				return fmt.Errorf("cannot make table: invalid row length found")
			}
			fmtRows[r] = make([]string, nc)
			for c, col := range row {
				s := fmt.Sprint(col)
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

	delete(m, "columns")

	for k, data := range m {
		// if reflect.TypeOf(data).Kind() != reflect.Slice {
		// 	return true, fmt.Errorf("%s not a slice", k)
		// }
		// rows := reflect.ValueOf(data)
		var rows []interface{}
		var ok bool
		if rows, ok = data.([]interface{}); !ok {
			return true, fmt.Errorf("%s not a slice", k)
		}

		v.addTable(columns, rows, false, true, w)
		break
	}

	// fmt.Print(columns)

	// if reflect.TypeOf(columns).Kind() != reflect.Slice {
	// 	return false, fmt.Errorf("columns not a slice")
	// }

	// slice := reflect.ValueOf(columns)
	// n := slice.Len()
	// if n == 0 {
	// 	return false, fmt.Errorf("zero comuns found")
	// }

	// var temp = make([]string, n)
	// for i := 0; i < n; i++ {
	// 	v := slice.Index(i).Interface()
	// 	if s, ok := v.(string); ok {
	// 		temp[i] = s
	// 	} else {
	// 		return false, fmt.Errorf("columns contains non string")
	// 	}
	// }

	// *lines = append(*lines, strings.Join(temp, ","))

	// delete(m, "columns")

	// for k, data := range m {
	// 	if reflect.TypeOf(data).Kind() != reflect.Slice {
	// 		return true, fmt.Errorf("%s not a slice", k)
	// 	}
	// 	rows := reflect.ValueOf(data)
	// 	nrows := rows.Len()
	// 	for r := 0; r < nrows; r++ {
	// 		row := rows.Index(r).Interface()

	// 		if reflect.TypeOf(row).Kind() != reflect.Slice {
	// 			return true, fmt.Errorf("row not a slice")
	// 		}
	// 		cols := reflect.ValueOf(row)

	// 		ncols := cols.Len()
	// 		if n != ncols {
	// 			return true, fmt.Errorf("number of columns does not equel values")
	// 		}
	// 		var temp = make([]string, n)
	// 		for i := 0; i < ncols; i++ {
	// 			temp[i] = escapeCsv(fmt.Sprint(cols.Index(i).Interface()))
	// 		}
	// 		*lines = append(*lines, strings.Join(temp, ","))
	// 	}
	// }
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
	}
	v.newLine()
	err := v.addPretty(w)
	if err != nil {
		v.error(err, w)
	}
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
