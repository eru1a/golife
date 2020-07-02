package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

var (
	deadColor        = tcell.NewRGBColor(50, 50, 50)
	aliveColor       = tcell.NewRGBColor(100, 200, 100)
	cursorDeadColor  = tcell.NewRGBColor(50, 50, 250)
	cursorAliveColor = tcell.NewRGBColor(50, 250, 250)
	strColor         = tcell.NewRGBColor(255, 255, 0)
)

type cursor struct {
	x, y int
}

type Game struct {
	screen     tcell.Screen
	board      *Board
	running    bool
	mu         sync.Mutex
	cursor     cursor
	pattern    *Pattern
	patterns   []*Pattern
	patternIdx int
}

var (
	cellPattern = &Pattern{
		Cells: &Cells{
			cells:  [][]bool{{true}},
			width:  1,
			height: 1,
		},
		name:    "cell",
		comment: []string{},
		credit:  "",
	}
)

func NewGame(s tcell.Screen, patterns []*Pattern) *Game {
	w, h := s.Size()
	w /= 2
	b := NewBoard(w, h)
	patterns = append([]*Pattern{cellPattern}, patterns...)
	return &Game{
		screen:     s,
		board:      b,
		running:    false,
		mu:         sync.Mutex{},
		cursor:     cursor{},
		pattern:    cellPattern,
		patterns:   patterns,
		patternIdx: 0,
	}
}

func (g *Game) display() {
	g.screen.Clear()

	g.displayBoard()

	g.displayPattern()

	g.displayString(0, g.board.Height()-8, "s: start/stop")
	g.displayString(0, g.board.Height()-7, "n: next")
	g.displayString(0, g.board.Height()-6, "r: random")
	g.displayString(0, g.board.Height()-5, "space: set")
	g.displayString(0, g.board.Height()-4, fmt.Sprintf("p: pattern [%s]", g.pattern.name))
	g.displayString(0, g.board.Height()-3, fmt.Sprintf("cursor: [%d, %d]", g.cursor.x, g.cursor.y))
	g.displayString(0, g.board.Height()-2, fmt.Sprintf("generation: %d", g.board.Generation()))
	g.displayString(0, g.board.Height()-1, "q: quit")

	g.screen.Show()
}

func (g *Game) displayBoard() {
	for y := 0; y < g.board.Height(); y++ {
		for x := 0; x < g.board.Width(); x++ {
			bg := deadColor
			if g.board.Get(x, y) {
				bg = aliveColor
			}
			style := tcell.StyleDefault.Background(bg)
			g.screen.SetContent(x*2, y, ' ', nil, style)
			g.screen.SetContent(x*2+1, y, ' ', nil, style)
		}
	}
}

func (g *Game) displayCursor(x, y int) {
	if !g.board.CheckRange(x, y) {
		return
	}
	bg := cursorDeadColor
	if g.board.Get(x, y) {
		bg = cursorAliveColor
	}
	style := tcell.StyleDefault.Background(bg)
	g.screen.SetContent(x*2, y, ' ', nil, style)
	g.screen.SetContent(x*2+1, y, ' ', nil, style)
}

func (g *Game) displayPattern() {
	for y := 0; y < g.pattern.Height(); y++ {
		for x := 0; x < g.pattern.Width(); x++ {
			if !g.pattern.Get(x, y) {
				continue
			}
			g.displayCursor(g.cursor.x+x-g.pattern.Width()/2, g.cursor.y+y-g.pattern.Height()/2)
		}
	}
}

func (g *Game) displayString(x, y int, str string) {
	for i, b := range str {
		bg := deadColor
		onAlive := g.board.Get((x+i)/2, y)
		onCursor := (x+i)/2 == g.cursor.x && y == g.cursor.y
		switch {
		case onAlive && onCursor:
			bg = cursorAliveColor
		case onAlive:
			bg = aliveColor
		case onCursor:
			bg = cursorDeadColor
		default:
			bg = deadColor
		}
		style := tcell.StyleDefault.Foreground(strColor).Background(bg)
		g.screen.SetContent(x+i, y, b, nil, style)
	}
}

func (g *Game) set(c bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.board.Set(g.cursor.x, g.cursor.y, c)
}

func (g *Game) setPattern() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.board.SetPattern(g.cursor.x-g.pattern.Width()/2, g.cursor.y-g.pattern.Height()/2, g.pattern)
}

func (g *Game) next() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.board.Next()
}

func (g *Game) random() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.board.Random()
}

func (g *Game) Loop() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if g.running {
				g.next()
				g.display()
			}
		}
	}()

	g.display()
	for {
		ev := g.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch {
			case ev.Rune() == 's':
				g.running = !g.running
			case ev.Rune() == 'c':
				g.board.Clear()
				g.running = false
				g.display()
			case ev.Rune() == 'r':
				g.random()
				g.display()
			case ev.Rune() == 'n':
				g.next()
				g.display()
			case ev.Rune() == 'p':
				g.patternIdx = (g.patternIdx + 1) % len(g.patterns)
				g.pattern = g.patterns[g.patternIdx]
				g.display()
			case ev.Rune() == ' ':
				g.setPattern()
				g.display()
			case ev.Key() == tcell.KeyRight, ev.Rune() == 'l':
				g.cursor.x++
				if g.cursor.x >= g.board.Width() {
					g.cursor.x = g.board.Width() - 1
				}
				g.display()
			case ev.Key() == tcell.KeyLeft, ev.Rune() == 'h':
				g.cursor.x--
				if g.cursor.x < 0 {
					g.cursor.x = 0
				}
				g.display()
			case ev.Key() == tcell.KeyDown, ev.Rune() == 'j':
				g.cursor.y++
				if g.cursor.y >= g.board.Height() {
					g.cursor.y = g.board.Height() - 1
				}
				g.display()
			case ev.Key() == tcell.KeyUp, ev.Rune() == 'k':
				g.cursor.y--
				if g.cursor.y < 0 {
					g.cursor.y = 0
				}
				g.display()
			case ev.Key() == tcell.KeyCtrlC, ev.Rune() == 'q':
				return
			}
		case *tcell.EventMouse:
			// urxvt上だと何かおかしい？
			x, y := ev.Position()
			x /= 2
			if !g.board.CheckRange(x, y) {
				continue
			}
			g.cursor.x, g.cursor.y = x, y
			if ev.Buttons() == 0 {
				g.display()
				continue
			}
			if ev.Buttons()&tcell.Button1 != 0 {
				g.setPattern()
			} else {
				g.set(false)
			}
			g.display()
		}
	}
}
