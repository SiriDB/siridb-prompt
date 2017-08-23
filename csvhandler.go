package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func toCsv(m map[string]interface{}) (string, error) {
	l := make([]string, 0)
	if err := getLines(m, &l); err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		return strings.Join(l, "\r\n"), nil
	}
	return strings.Join(l, "\n"), nil
}

func getLines(m map[string]interface{}, l *[]string) error {
	if err := tryShow(m, l); err != nil {
		return err
	}

	if err := tryCalc(m, l); err != nil {
		return err
	}

	if err := tryMsg(m, l); err != nil {
		return err
	}

	if err := tryCount(m, l); err != nil {
		return err
	}

	if columns := getListColumns(m); columns != nil {
		if err := tryList(m, l, columns); err != nil {
			return err
		}
	} else {
		if err := trySelect(m, l); err != nil {
			return err
		}
	}

	if err := tryTimeit(m, l); err != nil {
		return err
	}

	return nil
}

func escapeCsv(i interface{}) string {
	s := fmt.Sprint(i)
	if strings.ContainsRune(s, '"') {
		return fmt.Sprintf("\"%s\"", strings.Replace(s, `"`, `""`, -1))
	}
	if strings.ContainsRune(s, ',') {
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

func getCsvLine(items ...interface{}) string {
	arr := make([]string, len(items))
	for i, itm := range items {
		arr[i] = escapeCsv(itm)
	}
	return strings.Join(arr, ",")
}

func addTable(columns []interface{}, rows []interface{}, mustSort bool, lines *[]string) error {
	var csvrows []string
	nc := len(columns)
	if nc == 0 {
		return fmt.Errorf("cannot make csv table: at least one column is required")
	}
	*lines = append(*lines, getCsvLine(columns...))

	for _, rowIf := range rows {
		if row, ok := rowIf.([]interface{}); ok {
			if nc != len(row) {
				return fmt.Errorf("cannot make csv table: invalid row length found")
			}
			csvrows = append(csvrows, getCsvLine(row...))
		}
	}
	if mustSort {
		sort.Slice(csvrows, func(i, j int) bool {
			return strings.Compare(csvrows[i], csvrows[j]) < 0
		})
	}
	*lines = append(*lines, csvrows...)

	return nil
}

func tryShow(m map[string]interface{}, lines *[]string) error {
	var err error
	if data, ok := m["data"]; ok {
		var columns = []interface{}{"name", "value"}
		var rows []interface{}
		if items, ok := data.([]interface{}); ok {
			for _, item := range items {
				if itm, ok := item.(map[string]interface{}); ok {
					if name, ok := itm["name"]; ok {
						if value, ok := itm["value"]; ok {
							rows = append(rows, []interface{}{name, value})
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
			err = addTable(columns, rows, true, lines)
		}
	}
	return err
}

func tryCalc(m map[string]interface{}, lines *[]string) error {
	if data, ok := m["calc"]; ok {
		n, ok := data.(int)
		if ok {
			*lines = append(*lines, getCsvLine("Timestamp", n))
			*lines = append(*lines, getCsvLine("UTC", fmtTimestampUTC(n)))
			*lines = append(*lines, getCsvLine("Local", fmtTimestamp(n)))
		}
	}
	return nil
}

func tryMsg(m map[string]interface{}, lines *[]string) error {
	options := [4]string{
		"success_msg",
		"error_msg",
		"help",
		"motd"}
	for _, option := range options {
		if message, ok := m[option]; ok {
			msg, ok := message.(string)
			if ok {
				*lines = append(*lines, getCsvLine(option, msg))
			}
		}
	}
	return nil
}

func tryCount(m map[string]interface{}, lines *[]string) error {
	for k, data := range m {
		switch k {
		case "calc":
			// skip
		default:
			n, ok := data.(int)
			if ok {
				*lines = append(*lines, getCsvLine(k, n))
			}
		}
	}
	return nil
}

func tryList(m map[string]interface{}, lines *[]string, columns []string) error {
	var err error
	for k, data := range m {
		switch k {
		case "__timeit__", "columns":
			// skip
		default:
			var rows []interface{}
			var ok bool
			if rows, ok = data.([]interface{}); ok {
				cols := make([]interface{}, len(columns))
				for i, c := range columns {
					cols[i] = c
				}
				err = addTable(cols, rows, true, lines)
			}
		}
	}
	return err
}

func trySelect(m map[string]interface{}, lines *[]string) error {
	for series, data := range m {
		if points, ok := data.([]interface{}); ok {
			for _, pointIf := range points {
				if point, test := pointIf.([]interface{}); test && len(point) == 2 {
					*lines = append(*lines, getCsvLine(append([]interface{}{series}, point...)...))
				}
			}
		}
	}
	return nil
}

func tryTimeit(m map[string]interface{}, lines *[]string) error {
	var err error
	if data, ok := m["__timeit__"]; ok {
		var columns = []interface{}{"server", "time"}
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
			*lines = append(*lines, "")
			err = addTable(columns, rows, true, lines)
		}
	}
	return err
}

func parseCsv(r io.Reader) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	reader := csv.NewReader(r)

	record, err := reader.Read()
	if err == io.EOF {
		return nil, fmt.Errorf("no csv data found")
	}
	if err != nil {
		return nil, err
	}
	if record[0] == "" {
		err = readTable(&data, record, reader)
	} else if len(record) == 3 {
		err = readFlat(&data, record, reader)
	} else {
		err = fmt.Errorf("unknown csv layout received")
	}
	return data, err
}

func parseCsvVal(inp string) interface{} {
	if i, err := strconv.Atoi(inp); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(inp, 64); err == nil {
		return f
	}
	return inp
}

func readTable(data *map[string]interface{}, record []string, reader *csv.Reader) error {
	if len(record) < 2 {
		return fmt.Errorf("missing series in csv table")
	}

	arr := make([][][2]interface{}, len(record)-1)

	for n := 1; n < len(record); n++ {
		(*data)[record[n]] = &arr[n-1]
	}
	for n := 2; ; n++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ts, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return fmt.Errorf("expecting a time-stamp in column zero at line %d", n)
		}
		for i := 1; i < len(record); i++ {
			arr[i-1] = append(arr[i-1], [2]interface{}{ts, parseCsvVal(record[i])})
		}
	}
	return nil
}

func readFlat(data *map[string]interface{}, record []string, reader *csv.Reader) error {
	appendFlatRecord(data, record, 1)
	for n := 2; ; n++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := appendFlatRecord(data, record, n); err != nil {
			return err
		}

	}
	return nil
}

func appendFlatRecord(data *map[string]interface{}, record []string, n int) error {
	var points *[][2]interface{}
	p, ok := (*data)[record[0]]
	if ok {
		points = p.(*[][2]interface{})
	} else {
		newPoints := make([][2]interface{}, 0)
		(*data)[record[0]] = &newPoints
		points = &newPoints
	}
	ts, err := strconv.ParseUint(record[1], 10, 64)
	if err != nil {
		return fmt.Errorf("expecting a time-stamp in column one at line %d", n)
	}
	*points = append(*points, [2]interface{}{ts, parseCsvVal(record[2])})
	return nil
}
