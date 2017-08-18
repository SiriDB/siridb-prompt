package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/transceptor-technology/go-siridb-connector"
	"github.com/transceptor-technology/goleri"

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
	xHistory  = xApp.Flag("history", "Number of command in history. A value of 0 disables history.").Default("1000").Uint16()
	xTimeout  = xApp.Flag("timeout", "Query timeout in seconds.").Default("60").Uint16()
	xServers  = xApp.Flag("servers", "Server(s) to connect to. Multiple servers are allowed and should be separated with a comma. (syntax: --servers=host[:port]").Short('s').String()
	xVersion  = xApp.Flag("version", "Print version information and exit.").Short('v').Bool()
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

const coldef = termbox.ColorDefault
const cViewLog = 0
const cViewOutput = 1

var logger = newLogView()
var outv = newView()
var client *siridb.Client
var currentView = cViewLog
var outPrompt = newPrompt(">>> ", coldef|termbox.AttrBold, coldef)
var his *history
var siriGrammar = SiriGrammar()

func logHandle(logCh chan string) {
	for {
		msg := <-logCh
		logger.append(msg)
		draw()
	}
}

func draw() {

	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()
	x := 0

	var s, tmp string
	var fg termbox.Attribute

	if currentView == cViewLog {
		s = " Log (ESC / CTRL+L close log, CTRL+Q quit)"
	} else {
		s = " Output (CTRL+L view log, CTRL+J copy to clipboard, CTRL+Q quit)"
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
		logger.draw(w, h)
	} else {
		outv.draw(w, h)
		outPrompt.draw(0, h-1, w, coldef, coldef)
	}

	termbox.Flush()
}

func sendCommand() int {
	s := strings.TrimSpace(string(outPrompt.text))
	if strings.Compare(s, "exit") == 0 {
		return 1
	}
	his.insert(s)
	q := newQuery(s)
	q.parse(*xTimeout)
	w, _ := termbox.Size()
	outv.append(q, w)
	outPrompt.deleteAllRunes()

	return 0
}

func toClipboard() {
	var err error
	var s string

	if outv.query == nil {
		err = fmt.Errorf("nothing to copy")
	} else {
		s, err = outv.query.json()
		if err == nil {
			err = clipboard.WriteAll(s)
		}
	}
	if err == nil {
		logger.append(fmt.Sprintf("successfully copied last result to clipboard"))
	} else {
		logger.append(fmt.Sprintf("cannot copy to clipboard: %s", err.Error()))
	}
}

func getCompletions(p *prompt) []completion {
	q := p.textBeforeCursor()
	res, err := siriGrammar.Parse(q)
	if err != nil {
		logger.append(fmt.Sprintf("goleri parse error: %s", err.Error()))
		return nil
	}

	var completions []completion
	rest := q[res.Pos():]
	if strings.HasPrefix("exit", q) {
		completions = append(completions, completion{
			text:     "exit",
			display:  "exit",
			startPos: -len(q),
		})
	}

	for _, elem := range res.GetExpecting() {
		if kw, ok := elem.(*goleri.Keyword); ok {
			word := kw.GetKeyword()
			if len(rest) == 0 && len(q) > 0 && q[len(q)-1] == ' ' || len(rest) > 0 && strings.HasPrefix(word, rest) {
				completions = append(completions, completion{
					text:     fmt.Sprintf("%s ", word),
					display:  word,
					startPos: -len(rest),
				})
			}
		}
	}

	return completions
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
	servers, err = getServers("localhost:9000,localhost:9001")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if len(*xUser) == 0 {
		*xUser = "iris"
	}
	if len(*xPassword) == 0 {
		*xPassword = "siri"
	}
	if len(*xDbname) == 0 {
		*xDbname = "dbtest"
	}

	var historyFnP *string
	historyFn, err := homedir.Dir()
	if err == nil {
		historyFn = path.Join(historyFn, ".siridb-prompt", fmt.Sprintf("%s@%s.history.1", *xUser, *xDbname))
		historyFnP = &historyFn
	}

	his = newHistory(int(*xHistory), historyFnP)
	his.load()
	defer his.save()

	client = siridb.NewClient(
		*xUser,                      // user
		*xPassword,                  // password
		*xDbname,                    // database
		serversToInterface(servers), // siridb server(s)
		logCh, // optional log channel
	)

	client.Connect()
	if client.IsAvailable() {
		currentView = cViewOutput
	}
	outPrompt.completer = getCompletions

	defer client.Close()

	draw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if currentView == cViewLog {
				switch ev.Key {
				case termbox.KeyCtrlQ:
					break mainloop
				case termbox.KeyCtrlL, termbox.KeyEsc:
					currentView = cViewOutput
				}
			} else if currentView == cViewOutput {
				switch ev.Key {
				case termbox.KeyCtrlQ:
					break mainloop
				case termbox.KeyCtrlL:
					currentView = cViewLog
				case termbox.KeyTab:
					// auto completion
				case termbox.KeyCtrlJ:
					toClipboard()
				case termbox.KeyEnter:
					if sendCommand() == 1 {
						break mainloop
					}
				default:
					outPrompt.parse(ev)
				}
			}
		case termbox.EventMouse:
			if currentView == cViewLog {
				switch ev.Key {
				case termbox.MouseWheelUp:
					logger.up()
				case termbox.MouseWheelDown:
					logger.down()
				}
			} else if currentView == cViewOutput {
				switch ev.Key {
				case termbox.MouseWheelUp:
					outv.up()
				case termbox.MouseWheelDown:
					outv.down()
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw()
	}
}
