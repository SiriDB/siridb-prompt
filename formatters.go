package main

import (
	"fmt"
	"math"
	"time"
)

func getFormatter(s string) func(i interface{}) string {
	switch s {
	case "start", "end", "timestamp":
		return fmtTimestamp
	case "size", "buffer_size":
		return fmtSizeBytes
	case "received_points", "selected_points", "list_limit", "select_points_limit":
		return fmtLargeNum
	case "mem_usage":
		return fmtSizeMb
	case "uptime":
		return fmtDuration
	case "drop_threshold":
		return fmtPercentage
	case "time_precision":
		return fmtTimePrecision
	case "time":
		return fmtTimeit
	default:
		return func(i interface{}) string { return fmt.Sprint(i) }
	}
}

func fmtTimeit(i interface{}) string {
	f, ok := i.(float64)
	if ok {
		return fmt.Sprintf("%.6f seconds", f)

	}
	return fmt.Sprint(i)
}

func fmtTimePrecision(i interface{}) string {
	tp, ok := i.(string)
	if ok {
		switch tp {
		case "s":
			return fmt.Sprintf("%s (%s)", tp, "second")
		case "ms":
			return fmt.Sprintf("%s (%s)", tp, "millisecond")
		case "us":
			return fmt.Sprintf("%s (%s)", tp, "microsecond")
		case "ns":
			return fmt.Sprintf("%s (%s)", tp, "namesecond")
		}
	}
	return fmt.Sprint(i)
}

func fmtPercentage(i interface{}) string {
	p, ok := i.(float64)
	if ok {
		return fmt.Sprintf("%v (%d%%)", i, int(p*100))
	}
	return fmt.Sprint(i)
}

func fmtSizeBytes(i interface{}) string {
	size, ok := i.(int)
	if ok {
		lookup := "BKMGTPEZYXWVU"
		if size > 0 {
			i := int(min(int(math.Log(float64(size)))/int(math.Log(1024)), 12))
			size = int(size * 100 / int(math.Pow(1024, float64(i))) / 100)
			if i > 0 {
				return fmt.Sprintf("%d %cB", size, lookup[i])
			}
		}
		return fmt.Sprintf("%d bytes", size)
	}
	return fmt.Sprint(i)
}

func fmtSizeMb(i interface{}) string {
	size, ok := i.(int)
	if ok {
		size *= int(math.Pow(1024, 2))
		return fmtSizeBytes(size)
	}
	return fmt.Sprint(i)
}

func fmtTimestamp(i interface{}) string {
	ts, ok := i.(int)
	if ok && timePrecision != nil {
		switch *timePrecision {
		case "s":
			return time.Unix(int64(ts), 0).Format(time.UnixDate)
		case "ms":
			s := int64(ts / 1000)
			ns := int64(ts % 1000 * 1000000)
			return time.Unix(s, ns).Format(time.UnixDate)
		case "us":
			s := int64(ts / 1000000)
			ns := int64(ts % 1000000 * 1000)
			return time.Unix(s, ns).Format(time.UnixDate)
		case "ns":
			s := int64(ts / 1000000000)
			ns := int64(ts % 1000000000)
			return time.Unix(s, ns).Format(time.UnixDate)
		}
	}
	return fmt.Sprint(i)
}

func fmtTimestampUTC(i interface{}) string {
	ts, ok := i.(int)
	if ok && timePrecision != nil {
		switch *timePrecision {
		case "s":
			return time.Unix(int64(ts), 0).UTC().Format(time.RFC1123)
		case "ms":
			s := int64(ts / 1000)
			ns := int64(ts % 1000 * 1000000)
			return time.Unix(s, ns).UTC().Format(time.RFC1123)
		case "us":
			s := int64(ts / 1000000)
			ns := int64(ts % 1000000 * 1000)
			return time.Unix(s, ns).UTC().Format(time.RFC1123)
		case "ns":
			s := int64(ts / 1000000000)
			ns := int64(ts % 1000000000)
			return time.Unix(s, ns).UTC().Format(time.RFC1123)
		}
	}
	return fmt.Sprint(i)
}

func fmtLargeNum(i interface{}) string {
	s := []rune(fmt.Sprint(i))
	l := len(s)

	if _, ok := i.(int); ok && l > 6 {
		var res []rune
		for n, c := range s {
			res = append(res, c)
			if l-n > 1 && (l-n)%3 == 1 {
				res = append(res, '.')
			}
		}
		s = res
	}
	return string(s)
}

func fmtDuration(i interface{}) string {
	seconds, ok := i.(int)
	if ok {
		if seconds == 1 {
			return "1 second"
		}
		if seconds <= 60*2 {
			return fmt.Sprintf("%d seconds", seconds)
		}
		if seconds <= 3600*2 {
			return fmt.Sprintf("%d minutes", seconds/60)
		}
		if seconds <= 86400*2 {
			return fmt.Sprintf("%d hours", seconds/3600)
		}
		if seconds <= 2592000*2 {
			return fmt.Sprintf("%d days", seconds/86400)
		}
		if seconds <= 31557600*2 {
			return fmt.Sprintf("%d months", seconds/2592000)
		}
		return fmt.Sprintf("%d years", seconds/31557600)
	}
	return fmt.Sprint(i)
}
