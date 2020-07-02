package main

import (
	"math/rand"
)

type Board struct {
	*Cells
	generation int
}

func NewBoard(width, height int) *Board {
	return &Board{NewCells(width, height), 0}
}

func (b *Board) Generation() int {
	return b.generation
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
			if !b.CheckRange(nx, ny) {
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
				next.cells[y][x] = b.cells[y][x]
			case 3:
				next.cells[y][x] = true
			default:
				next.cells[y][x] = false
			}
		}
	}
	b.cells = next.cells
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

func (b *Board) SetPattern(x, y int, p *Pattern) {
	for y2 := 0; y2 < p.Height(); y2++ {
		for x2 := 0; x2 < p.Width(); x2++ {
			if !b.CheckRange(x+x2, y+y2) {
				// panic(fmt.Sprint(x+x2, y+y2))
				continue
			}
			if !p.Get(x2, y2) {
				continue
			}
			b.Set(x+x2, y+y2, true)
		}
	}
}
