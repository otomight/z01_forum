package utils

func ToCSV(elems []string) string {
	var	i		int
	var	result	string

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
