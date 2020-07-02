package main

type Cells struct {
	cells  [][]bool
	width  int
	height int
}

func NewCells(width, height int) *Cells {
	var cells [][]bool = make([][]bool, height)
	for i := 0; i < height; i++ {
		cells[i] = make([]bool, width)
	}
	return &Cells{
		cells:  cells,
		width:  width,
		height: height,
	}
}

func (c *Cells) Width() int {
	return c.width
}

func (c *Cells) Height() int {
	return c.height
}

func (c *Cells) Set(x, y int, cell bool) {
	c.cells[y][x] = cell
}

func (c *Cells) Get(x, y int) bool {
	return c.cells[y][x]
}

func (c *Cells) CheckRange(x, y int) bool {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return false
	}
	return true
}
