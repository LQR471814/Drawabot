package lib

func max(numbers ...int) int {
	largest := numbers[0]
	for n := range numbers {
		if n > largest {
			largest = n
		}
	}

	return largest
}
