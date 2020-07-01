package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

func display(b *Board) {
	for y := 0; y < b.Height(); y++ {
		for x := 0; x < b.Width(); x++ {
			if b.Get(x, y) {
				fmt.Print("o")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	defer s.Fini()

	if err := s.Init(); err != nil {
		panic(err)
	}
	s.EnableMouse()

	g := NewGame(s)
	g.Loop()
}
