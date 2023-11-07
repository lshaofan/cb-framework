package slice

import "reflect"

// DiffSlice 查找切片中不同的元素 此方法性能较差,用于切片元素类型为interface{}的情况
func DiffSlice[T any](slice1, slice2 []T) []T {
	result := make([]T, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if reflect.DeepEqual(v1, v2) {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}

// DiffUintSlice 查找uint切片中不同的元素
func DiffUintSlice(slice1, slice2 []uint) []uint {
	result := make([]uint, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if v1 == v2 {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}

// DiffIntSlice 查找int切片中不同的元素
func DiffIntSlice(slice1, slice2 []int) []int {
	result := make([]int, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if v1 == v2 {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}

// DiffStringSlice 查找string切片中不同的元素
func DiffStringSlice(slice1, slice2 []string) []string {
	result := make([]string, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if v1 == v2 {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}

// DiffInt64Slice 查找int64切片中不同的元素
func DiffInt64Slice(slice1, slice2 []int64) []int64 {
	result := make([]int64, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if v1 == v2 {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}

// DiffInt32Slice 查找int32切片中不同的元素
func DiffInt32Slice(slice1, slice2 []int32) []int32 {
	result := make([]int32, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if v1 == v2 {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}

// DiffInt16Slice 查找int16切片中不同的元素
func DiffInt16Slice(slice1, slice2 []int16) []int16 {
	result := make([]int16, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if v1 == v2 {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}

// DiffInt8Slice 查找int8切片中不同的元素
func DiffInt8Slice(slice1, slice2 []int8) []int8 {
	result := make([]int8, 0)
	for _, v1 := range slice1 {
		exists := false
		for _, v2 := range slice2 {
			// 判断值是否相等
			if v1 == v2 {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, v1)
		}
	}
	return result
}
