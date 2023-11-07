package slice

// RemoveRepBySlice 去重切片中的元素方法
func RemoveRepBySlice[T any](slc []T) []T {

	m := make(map[any]bool)
	for _, v := range slc {
		m[v] = true
	}
	result := make([]T, 0, len(m))
	for k := range m {
		result = append(result, k.(T))
	}
	return result
}

// RemoveRepByUintSlice 去重uint切片中的元素方法
func RemoveRepByUintSlice(slc []uint) []uint {
	m := make(map[uint]bool)
	for _, v := range slc {
		m[v] = true
	}
	result := make([]uint, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result

}

// RemoveRepByIntSlice 去重int切片中的元素方法
func RemoveRepByIntSlice(slc []int) []int {
	m := make(map[int]bool)
	for _, v := range slc {
		m[v] = true
	}
	result := make([]int, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// RemoveRepByStringSlice 去重string切片中的元素方法
func RemoveRepByStringSlice(slc []string) []string {
	m := make(map[string]bool)
	for _, v := range slc {
		m[v] = true
	}
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
