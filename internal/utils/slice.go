package utils

func ToCSV(elems []string) string {
	var i int
	var result string

	i = 0
	for i < len(elems) {
		result += elems[i]
		i++
		if i != len(elems) {
			result += ", "
		}
	}
	return result
}

func ToMatrix(elems []string) [][]any {
	var result [][]any
	var i int

	result = make([][]any, len(elems))
	for i = 0; i < len(elems); i++ {
		result[i] = []any{elems[i]}
	}
	return result
}
