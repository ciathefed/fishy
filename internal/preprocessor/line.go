package preprocessor

import "fmt"

type Line struct {
	data string
	n    int
}

func (l Line) Print() {
	fmt.Printf("%d | %s\n", l.n, l.data)
}
