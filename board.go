package main

import (
	"fmt"
	"math/rand"
)

type Cells [][]bool

func NewCells(width, height int) Cells {
	var cells Cells = make([][]bool, height)
	for i := 0; i < height; i++ {
		cells[i] = make([]bool, width)
	}
	return cells
}

type Board struct {
	cells      Cells
	width      int
	height     int
	generation int
}

func NewBoard(width, height int) *Board {
	return &Board{
		cells:      NewCells(width, height),
		width:      width,
		height:     height,
		generation: 0,
	}
}

func (b *Board) Width() int {
	return b.width
}

func (b *Board) Height() int {
	return b.height
}

func (b *Board) Generation() int {
	return b.generation
}

func (b *Board) checkRange(x, y int) bool {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return false
	}
	return true
}

func (b *Board) Set(x, y int, c bool) {
	b.cells[y][x] = c
}

func (b *Board) Get(x, y int) bool {
	if !b.checkRange(x, y) {
		panic(fmt.Sprint(x, b.width, y, b.height))
	}
	return b.cells[y][x]
}

func (b *Board) Toggle(x, y int) {
	b.cells[y][x] = !b.cells[y][x]
}

func (b *Board) Around(x, y int) int {
	cnt := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx := x + dx
			ny := y + dy
			if !b.checkRange(nx, ny) {
				continue
			}
			if b.Get(nx, ny) {
				cnt++
			}
		}
	}
	return cnt
}

func (b *Board) Next() {
	next := NewCells(b.width, b.height)
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			switch b.Around(x, y) {
			case 2:
				next[y][x] = b.cells[y][x]
			case 3:
				next[y][x] = true
			default:
				next[y][x] = false
			}
		}
	}
	b.cells = next
	b.generation++
}

func (b *Board) Random() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if rand.Int()%10 == 0 {
				b.cells[y][x] = true
			}
		}
	}
}

func (b *Board) Clear() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x] = false
		}
	}
	b.generation = 0
}
