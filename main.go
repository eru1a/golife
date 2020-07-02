package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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

	patterns := []*Pattern{}
	filepath.Walk("./patterns", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pattern, err := NewPatternFromRLE(f)
		if err != nil {
			panic(err)
		}
		patterns = append(patterns, pattern)
		return nil
	})
	g := NewGame(s, patterns)
	g.Loop()
}
