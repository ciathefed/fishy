package preprocessor

type Macro struct {
	line   int
	col    int
	params []string
	lines  []string
}
