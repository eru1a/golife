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
	screen  tcell.Screen
	board   *Board
	running bool
	mu      sync.Mutex
	cursor  cursor
}

func NewGame(s tcell.Screen) *Game {
	w, h := s.Size()
	w /= 2
	b := NewBoard(w, h)
	return &Game{
		screen:  s,
		board:   b,
		running: false,
		mu:      sync.Mutex{},
		cursor:  cursor{},
	}
}

func (g *Game) display() {
	g.screen.Clear()

	g.displayBoard()

	g.displayCursor()

	g.displayString(0, g.board.Height()-6, "s: start/stop")
	g.displayString(0, g.board.Height()-5, "n: next")
	g.displayString(0, g.board.Height()-4, "r: random")
	g.displayString(0, g.board.Height()-3, "space: toggle")
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

func (g *Game) displayCursor() {
	x, y := g.cursor.x, g.cursor.y
	bg := cursorDeadColor
	if g.board.Get(x, y) {
		bg = cursorAliveColor
	}
	style := tcell.StyleDefault.Background(bg)
	g.screen.SetContent(x*2, y, ' ', nil, style)
	g.screen.SetContent(x*2+1, y, ' ', nil, style)
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

func (g *Game) Loop() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if g.running {
				g.mu.Lock()
				g.board.Next()
				g.mu.Unlock()
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
				g.mu.Lock()
				g.board.Random()
				g.mu.Unlock()
				g.display()
			case ev.Rune() == 'n':
				g.mu.Lock()
				g.board.Next()
				g.mu.Unlock()
				g.display()
			case ev.Rune() == ' ':
				g.mu.Lock()
				g.board.Toggle(g.cursor.x, g.cursor.y)
				g.mu.Unlock()
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
			if !g.board.checkRange(x, y) {
				continue
			}
			g.cursor.x, g.cursor.y = x, y
			if ev.Buttons() == 0 {
				g.display()
				continue
			}
			cell := false
			if ev.Buttons()&tcell.Button1 != 0 {
				cell = true
			}
			g.mu.Lock()
			g.board.Set(x, y, cell)
			g.mu.Unlock()
			g.display()
		}
	}
}
