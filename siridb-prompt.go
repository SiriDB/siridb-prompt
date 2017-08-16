package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/transceptor-technology/go-siridb-connector"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// AppVersion exposes version information
const AppVersion = "2.1.0"

var (
	xApp      = kingpin.New("siridb-admin", "Tool for communicating with a SiriDB database.")
	xUser     = xApp.Flag("user", "Database user.").Short('u').String()
	xPassword = xApp.Flag("password", "Password for the database user.").Short('p').String()
	xDbname   = xApp.Flag("dbname", "Database name.").Short('d').String()
	xServers  = xApp.Flag("servers", "Server(s) to connect to. Multiple servers are allowed and should be separated with a comma. (syntax: --servers=host[:port]").Short('s').String()
	xVersion  = xApp.Flag("version", "Print version information and exit.").Short('v').Bool()
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

type server struct {
	host string
	port uint16
}

type logEntry struct {
	dtime time.Time
	msg   string
}

type logView struct {
	pos     int
	entries []logEntry
}

var logger = logView{pos: -1, entries: make([]logEntry, 0)}

var client *siridb.Client

const cViewLog = 0
const cViewOutput = 1

var currentView = cViewLog
var X = 8

func getHostAndPort(addr string) (server, error) {
	parts := strings.Split(addr, ":")
	// IPv4
	if len(parts) == 1 {
		return server{parts[0], 9000}, nil
	}
	if len(parts) == 2 {
		u, err := strconv.ParseUint(parts[1], 10, 16)
		return server{parts[0], uint16(u)}, err
	}
	// IPv6
	if addr[0] != '[' {
		return server{fmt.Sprintf("[%s]", addr), 9000}, nil
	}
	if addr[len(addr)-1] == ']' {
		return server{addr, 9000}, nil
	}
	u, err := strconv.ParseUint(parts[len(parts)-1], 10, 16)
	addr = strings.Join(parts[:len(parts)-1], ":")

	return server{addr, uint16(u)}, err
}

func getServers(addrstr string) ([]server, error) {
	arr := strings.Split(addrstr, ",")
	servers := make([]server, len(arr))
	for i, addr := range arr {
		addr = strings.TrimSpace(addr)
		server, err := getHostAndPort(addr)
		if err != nil {
			return nil, err
		}
		servers[i] = server
	}
	return servers, nil
}

func serversToInterface(servers []server) [][]interface{} {
	ret := make([][]interface{}, len(servers))
	for i, svr := range servers {
		ret[i] = make([]interface{}, 2)
		ret[i][0] = svr.host
		ret[i][1] = int(svr.port)
	}
	return ret
}

func quit(err error) {
	rc := 0
	if err != nil {
		fmt.Printf("%s\n", err)
		rc = 1
	}

	if client != nil {
		client.Close()
	}

	os.Exit(rc)
}

func logHandle(logCh chan string) {
	for {
		msg := <-logCh
		logger.entries = append(logger.entries, logEntry{time.Now(), msg})
		draw()
	}
}

func draw() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()
	x := 0

	var s, tmp string
	var fg termbox.Attribute

	if currentView == cViewLog {
		s = " Connection Log (ESC or CTRL+L to close)"
	} else {
		s = " Output (CTRL+L to open connection log, CTRL+C to quit)"
	}
	for _, c := range s {
		termbox.SetCell(x, 0, c, termbox.ColorBlack, termbox.ColorWhite)
		x++
	}
	s = fmt.Sprintf(" <%s@%s> status: ", *xUser, *xDbname)
	if client.IsAvailable() {
		tmp = "OK "
		fg = termbox.ColorGreen
	} else {
		tmp = "NO CONNECTION "
		fg = termbox.ColorRed
	}
	end := w - len(s) - len(tmp)
	for ; x < end; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorBlack, termbox.ColorWhite)
	}

	for _, c := range s {
		termbox.SetCell(x, 0, c, termbox.ColorBlack, termbox.ColorWhite)
		x++
	}

	for _, c := range tmp {
		termbox.SetCell(x, 0, c, fg, termbox.ColorWhite)
		x++
	}

	if currentView == cViewLog {
		y := 1
		pos := logger.pos
		if pos == -1 {
			pos = len(logger.entries)
		}
		start := pos - (h - 1)

		if start < 0 {
			start = 0
		}

		for _, entry := range logger.entries[start:pos] {
			x = 0
			s = fmt.Sprintf("%s %s", entry.dtime.Format(time.UnixDate), entry.msg)
			for _, c := range s {
				termbox.SetCell(x, y, c, coldef, coldef)
				x++
			}
			y++

		}
	}

	termbox.Flush()
}

func main() {
	rand.Seed(time.Now().Unix())

	_, err := xApp.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if *xVersion {
		fmt.Printf("Version: %s\n", AppVersion)
		os.Exit(0)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	logCh := make(chan string)

	go logHandle(logCh)

	var servers []server
	servers, err = getServers("localhost:9000")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	*xUser = "iris"
	*xPassword = "siri"
	*xDbname = "dbtest"

	client = siridb.NewClient(
		*xUser,                      // user
		*xPassword,                  // password
		*xDbname,                    // database
		serversToInterface(servers), // siridb server(s)
		logCh, // optional log channel
	)

	client.Connect()

	defer client.Close()

	draw()

	const coldef = termbox.ColorDefault

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop

			default:
				if ev.Ch != 0 {
					fmt.Println(ev.Ch)
					termbox.Clear(coldef, coldef)
					termbox.SetCell(X, 8, '!', coldef, coldef)
					termbox.Flush()
				}
			}
		case termbox.EventMouse:
			switch ev.Key {
			case termbox.MouseWheelUp:
				if currentView == cViewLog {
					if logger.pos > 0 {
						logger.pos--
					} else if logger.pos == -1 {
						logger.pos = len(logger.entries) - 1
					}
				}
			case termbox.MouseWheelDown:
				if currentView == cViewLog {
					if logger.pos != -1 {
						logger.pos++
					}
					if logger.pos >= len(logger.entries) {
						logger.pos = -1
					}
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw()
	}
}
