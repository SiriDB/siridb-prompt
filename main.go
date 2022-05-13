package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/SiriDB/go-siridb-connector"
	"github.com/atotto/clipboard"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/transceptor-technology/goleri"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// AppVersion exposes version information
const AppVersion = "2.1.11"

var (
	xApp      = kingpin.New("siridb-admin", "Tool for communicating with a SiriDB database.")
	xDbname   = xApp.Flag("dbname", "Database name.").Short('d').Required().String()
	xServers  = xApp.Flag("servers", "Server(s) to connect to. Multiple servers are allowed and should be separated with a comma. (syntax: --servers=host[:port]").Short('s').Required().String()
	xUser     = xApp.Flag("user", "Database user.").Short('u').Required().String()
	xPassword = xApp.Flag("password", "Password for the database user.").Short('p').String()
	xRun      = xApp.Flag("run", "Run single command.").Short('r').String()
	xHistory  = xApp.Flag("history", "Number of command in history. A value of 0 disables history.").Default("1000").Uint16()
	xTimeout  = xApp.Flag("timeout", "Query timeout in seconds.").Default("60").Uint16()
	xJSON     = xApp.Flag("json", "Raw JSON output.").Bool()
	xMouse    = xApp.Flag("mouse", "Enable mouse support.").Bool()
	xVersion  = xApp.Flag("version", "Print version information and exit.").Short('v').Bool()
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

const coldef = termbox.ColorDefault
const colsel = termbox.ColorMagenta

const cViewLog = 0
const cViewOutput = 1

var logger = newLogView()
var isDrawwing = false
var outv = newView()
var mouseSelect = newMselect()
var client *siridb.Client
var currentView = cViewOutput
var outPrompt = newPrompt(">>> ", coldef|termbox.AttrBold, coldef)
var his *history
var siriGrammar = SiriGrammar()
var timePrecision *string

func drawPassPrompt(p *prompt) {
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()
	p.draw(0, 0, w, h, coldef, coldef)
	termbox.Flush()
}

func draw() {
	if isDrawwing {
		return
	}
	isDrawwing = true
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()
	x := 0

	var s, tmp string
	var fg termbox.Attribute

	if currentView == cViewLog {
		s = " Log (ESC / CTRL+L close log, CTRL+Q/D quit)"
	} else {
		s = " Output (CTRL+L view log, CTRL+J / CTRL+C copy to clipboard, CTRL+Q/D quit)"
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
		mouseSelect.draw(w, h)
		outPrompt.draw(0, h-1, w, h, coldef, coldef)
	}

	termbox.Flush()
	isDrawwing = false
}

func sendCommand() int {
	s := strings.TrimSpace(string(outPrompt.text))
	if strings.Compare(s, "exit") == 0 {
		return 1
	}
	his.insert(s)
	q := newQuery(s)
	if strings.HasPrefix(q.req, "dump ") {
		fn := strings.TrimSpace(q.req[4:])
		if outv.query != nil {
			q.res = outv.query.res
		}
		q.err = q.dumpToFile(fn)
	} else {
		q.parse(*xTimeout)
	}
	w, _ := termbox.Size()
	outv.append(q, w)
	outPrompt.deleteAllRunes()

	return 0
}

func toClipboard(to string) {
	var err error
	var s string

	if outv.query == nil {
		err = fmt.Errorf("nothing to copy")
	} else {
		switch to {
		case "JSON":
			s, err = outv.query.json()
		case "CSV":
			s, err = outv.query.csv()
		case "SELECTION":
			s = string(mouseSelect.getSelection())
		default:
			panic(fmt.Errorf("unexpected format: %s", to))
		}

		if err == nil {
			err = clipboard.WriteAll(s)
		}
	}
	if err == nil {
		logger.ch <- fmt.Sprintf("successfully copied last result as %s to clipboard", to)
	} else {
		logger.ch <- fmt.Sprintf("cannot copy to clipboard: %s", err.Error())
	}
}

func getCompletions(p *prompt) []*completion {
	q := p.textBeforeCursor()
	res, err := siriGrammar.Parse(q)
	if err != nil {
		logger.ch <- fmt.Sprintf("goleri parse error: %s", err.Error())
		return nil
	}

	var completions []*completion
	rest := q[res.Pos():]
	trimmed := strings.TrimLeft(q, " ")
	if len(trimmed) < 4 && strings.HasPrefix("exit", trimmed) {
		compl := completion{
			text:     "exit",
			display:  "exit",
			startPos: len(trimmed),
		}
		completions = append(completions, &compl)
	}

	for _, prefix := range []string{"import", "dump"} {
		pps := fmt.Sprintf("%s ", prefix)
		if len(trimmed) < len(prefix) && strings.HasPrefix(prefix, trimmed) {
			compl := completion{
				text:     pps,
				display:  prefix,
				startPos: len(trimmed),
			}
			completions = append(completions, &compl)
		}

		if strings.HasPrefix(trimmed, pps) {
			var fn string
			p := strings.TrimLeft(trimmed[len(prefix):], " ")
			p, fn = path.Split(p)

			if len(p) == 0 {
				p = "."
			}

			if files, err := ioutil.ReadDir(p); err == nil {
				s := strings.TrimSpace(path.Join("..", " "))
				if strings.HasPrefix(s, fn) {
					compl := completion{
						text:     s,
						display:  s,
						startPos: len(fn),
					}
					completions = append(completions, &compl)
				}
				for _, f := range files {
					var s, d string
					if f.IsDir() {
						d = strings.TrimSpace(path.Join(f.Name(), " "))
						s = d[:len(d)-1]
					} else {
						d = f.Name()
						locase := strings.ToLower(d)
						if !strings.HasSuffix(locase, ".json") && !strings.HasSuffix(locase, ".csv") {
							continue
						}
						s = d
					}
					if !strings.HasPrefix(s, fn) {
						continue
					}

					compl := completion{
						text:     s,
						display:  d,
						startPos: len(fn),
					}
					completions = append(completions, &compl)
				}
			}
		}
	}

	for _, elem := range res.GetExpecting() {
		if kw, ok := elem.(*goleri.Keyword); ok {
			word := kw.GetKeyword()
			if len(p.text) == 0 || len(rest) == 0 && len(q) > 0 && q[len(q)-1] == ' ' || len(rest) > 0 && strings.HasPrefix(word, rest) {
				compl := completion{
					text:     fmt.Sprintf("%s ", word),
					display:  word,
					startPos: len(rest),
				}
				completions = append(completions, &compl)
			}
		}
	}
	return completions
}

func initConnect() int {
	var tp string

	res, err := client.Query("show time_precision", 10)
	if err != nil {
		logger.ch <- fmt.Sprintf("error reading time_precision: %s", err.Error())
		return 1
	}
	v, ok := res.(map[string]interface{})
	if !ok {
		logger.ch <- "error reading time_precision: missing 'map' in data"
		return 1
	}

	arr, ok := v["data"].([]interface{})
	if !ok || len(arr) != 1 {
		logger.ch <- "error reading time_precision: missing array 'data' or length 1 in map"
		return 1
	}

	tp, ok = arr[0].(map[string]interface{})["value"].(string)

	if !ok {
		logger.ch <- "error reading time_precision: cannot find time_precision in data"
		return 1
	}

	logger.ch <- fmt.Sprintf("finished reading time precision: '%s'", tp)
	timePrecision = &tp
	return 0
}

func main() {
	exitCode := 0

	defer func() {
		os.Exit(exitCode)
	}()

	rand.Seed(time.Now().Unix())
	var historyFn string
	var err error
	var servers []server
	var historyFnP *string

	_, err = xApp.Parse(os.Args[1:])
	oneCommandOnly := len(*xRun) > 0
	askPassword := len(*xPassword) == 0

	logger.setMode("CONSOLE")

	go logger.handle()

	if *xVersion {
		fmt.Printf("Version: %s\n", AppVersion)
		goto finished
	}

	if err != nil {
		fmt.Printf("%s\n", err)
		goto stoperr
	}

	if *xJSON {
		outv.setModeJSON()
	}

	servers, err = getServers(*xServers)
	if err != nil {
		fmt.Printf("error reading servers: %s\n", err.Error())
		goto stoperr
	}

	if askPassword {
		err = termbox.Init()
		if err != nil {
			panic(err)
		}
		termbox.SetInputMode(termbox.InputEsc)
		termbox.SetOutputMode(termbox.Output256)
		pp := newPrompt("Password: ", coldef|termbox.AttrBold, coldef)
		pp.hideText = true

	passloop:
		for {
			drawPassPrompt(pp)
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyCtrlC, termbox.KeyCtrlQ:
					termbox.Close()
					goto stoperr
				case termbox.KeyEnter:
					*xPassword = string(pp.text)
					break passloop
				default:
					pp.parse(ev)
				}
			case termbox.EventError:
				panic(ev.Err)
			}

		}
		termbox.Close()
	}

	if !oneCommandOnly {
		outPrompt.completer = getCompletions

		historyFn, err = homedir.Dir()
		if err == nil {
			historyFn = path.Join(historyFn, ".siridb-prompt", fmt.Sprintf("%s@%s.history.1", *xUser, *xDbname))
			historyFnP = &historyFn
		}

		his = newHistory(int(*xHistory), historyFnP)
		his.load()
		defer his.save()
	}

	client = siridb.NewClient(
		*xUser,                      // user
		*xPassword,                  // password
		*xDbname,                    // database
		serversToInterface(servers), // siridb server(s)
		logger.ch,                   // optional log channel
	)

	logger.ch <- fmt.Sprintf("connecting to database %s...", *xDbname)

	client.Connect()

	defer client.Close()

	if !client.IsAvailable() || initConnect() != 0 {
		goto stoperr
	}

	if oneCommandOnly {
		q := newQuery(*xRun)
		q.parse(*xTimeout)
		outv.append(q, 99999)

		if q.err == nil {
			if len(outv.lines) > 1 {
				for _, line := range outv.lines[1:] {
					fmt.Println(line)
				}
			}
			goto finished
		}
		goto stoperr
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	defer func() {
		termbox.Close()
		logger.setMode("CONSOLE")
	}()

	if *xMouse {
		termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	} else {
		termbox.SetInputMode(termbox.InputEsc)
	}

	termbox.SetOutputMode(termbox.Output256)
	logger.setMode("TERMBOX")

mainloop:
	for {
		draw()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if currentView == cViewLog {
				switch ev.Key {
				case termbox.KeyCtrlQ:
				case termbox.KeyCtrlD:
					break mainloop
				case termbox.KeyCtrlL, termbox.KeyEsc:
					currentView = cViewOutput
				case termbox.KeyArrowUp:
					logger.up()
				case termbox.KeyArrowDown:
					logger.down()
				case termbox.KeyPgdn:
					logger.pageDown()
				case termbox.KeyPgup:
					logger.pageUp()
				}
			} else if currentView == cViewOutput {
				switch ev.Key {
				case termbox.KeyCtrlQ:
				case termbox.KeyCtrlD:
					break mainloop
				case termbox.KeyCtrlL:
					mouseSelect.clear()
					currentView = cViewLog
				case termbox.KeyCtrlJ:
					toClipboard("JSON")
				case termbox.KeyCtrlC:
					if mouseSelect.hasSelection {
						toClipboard("SELECTION")
					} else {
						toClipboard("CSV")
					}
				case termbox.KeyCtrlV:
					if mouseSelect.hasSelection {
						outPrompt.hidePopup()
						for _, r := range mouseSelect.getSelection() {
							outPrompt.insertRune(r)
						}
					}
				case termbox.KeyEnter:
					mouseSelect.clear()
					outPrompt.clearCompletions()
					if sendCommand() == 1 {
						break mainloop
					}
				case termbox.KeyPgdn:
					mouseSelect.clear()
					outv.pageDown()
				case termbox.KeyPgup:
					mouseSelect.clear()
					outv.pageUp()
				case termbox.KeyArrowUp:
					if outPrompt.hasCompletions() {
						outPrompt.parse(ev)
					} else {
						outPrompt.setText(his.prev())
						outPrompt.clearCompletions()
					}
				case termbox.KeyArrowDown:
					if outPrompt.hasCompletions() {
						outPrompt.parse(ev)
					} else {
						outPrompt.setText(his.next())
						outPrompt.clearCompletions()
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
					mouseSelect.clear()
					outv.up()
				case termbox.MouseWheelDown:
					mouseSelect.clear()
					outv.down()
				case termbox.MouseLeft:
					outPrompt.hidePopup()
					mouseSelect.start(ev.MouseX, ev.MouseY)
				case termbox.MouseRelease:
					mouseSelect.end(ev.MouseX, ev.MouseY)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
	goto finished

stoperr:
	logger.toStdErr()
	exitCode = 1

finished:
}
