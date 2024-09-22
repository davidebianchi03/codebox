package utils

func ItemInIntegersList(list []int, item int) bool {
	for _, val := range list {
		if val == item {
			return true
		}
	}

	return false
}
