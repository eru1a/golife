package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Pattern struct {
	*Cells
	name    string
	comment []string
	credit  string
}

func NewPatternFromRLE(r io.Reader) (*Pattern, error) {
	rb := bufio.NewReader(r)
	var w, h int
	name := "unknown"
	comment := []string{}
	credit := ""
	for {
		line, err := rb.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSuffix(line, "\n")
		if strings.HasPrefix(line, "#") {
			switch {
			case strings.HasPrefix(line, "#N"):
				line = strings.TrimLeft(line, "#N")
				line = strings.TrimSpace(line)
				name = line
			case strings.HasPrefix(line, "#C"), strings.HasPrefix(line, "#c"):
				line = strings.TrimLeft(line, "#C")
				line = strings.TrimLeft(line, "#c")
				line = strings.TrimSpace(line)
				comment = append(comment, line)
			case strings.HasPrefix(line, "#O"):
				line = strings.TrimLeft(line, "#O")
				line = strings.TrimSpace(line)
				credit = line
			}
			continue
		} else if strings.HasPrefix(line, "x") {
			fmt.Sscanf(line, "x = %d, y = %d", &w, &h)
			break
		} else {
			return nil, fmt.Errorf("should be '#' or 'x': %s", line)
		}
	}
	cells := NewCells(w, h)
	var x, y, cnt int
	for {
		b, err := rb.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if b == '\n' {
			continue
		}
		if b == '\r' {
			continue
		}

		if '0' <= b && b <= '9' {
			cnt = cnt*10 + int(b-'0')
			continue
		}

		n := 1
		if cnt > 0 {
			n = cnt
			cnt = 0
		}

		switch b {
		case 'b':
			for i := 0; i < n; i++ {
				cells.Set(x, y, false)
				x++
			}
		case 'o':
			for i := 0; i < n; i++ {
				cells.Set(x, y, true)
				x++
			}
		case '$':
			x = 0
			for i := 0; i < n; i++ {
				y++
			}
		case '!':
		default:
			return nil, fmt.Errorf("unexpected: %s", string(b))
		}
	}
	return &Pattern{cells, name, comment, credit}, nil
}
