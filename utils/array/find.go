package array

// Find 根据条件在数组中查找元素
func Find[T any](arr []T, fn func(T) bool) T {
	for _, v := range arr {
		if fn(v) {
			return v
		}
	}
	var ret T
	return ret
}
