package framework

import (
	"github.com/lshaofan/cb-framework/utils/slice"
	"reflect"
	"testing"
)

func TestRemoveRepByUintSlice(t *testing.T) {
	// Positive test cases
	slc1 := []uint{1, 2, 3, 4, 5, 5}
	expected1 := []uint{1, 2, 3, 4, 5}
	result1 := slice.RemoveRepByUintSlice(slc1)
	// 判断是否还有重复元素
	if len(result1) != len(expected1) {
		t.Error("RemoveRepByUintSlice failed")
	}
}

// RemoveRepBySlice 测试
func TestRemoveRepBySlice(t *testing.T) {
	slc1 := []string{"1", "2", "3", "4", "5", "5"}
	expected1 := []string{"1", "2", "3", "4", "5"}
	result1 := slice.RemoveRepBySlice(slc1)
	// 判断是否还有重复元素
	if len(result1) != len(expected1) {
		t.Error("RemoveRepBySlice failed")
	}

	// 测试interface{}类型
	slc2 := []interface{}{"1", "2", "3", "4", "5", "5", 6, 6}
	expected2 := []interface{}{"1", "2", "3", "4", "5", 6}
	result2 := slice.RemoveRepBySlice(slc2)
	// 判断是否还有重复元素
	if len(result2) != len(expected2) {
		t.Error("RemoveRepBySlice failed")
	}
}

func TestDiffUintSlice(t *testing.T) {
	slice1 := []uint{1, 2, 3, 4, 5}
	slice2 := []uint{3, 4, 5, 6, 7}
	expected := []uint{1, 2}
	result := slice.DiffUintSlice(slice1, slice2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DiffIntSlice(%v, %v) = %v; want %v", slice1, slice2, result, expected)
	}
}

func TestDiffSlice(t *testing.T) {
	// 正向测试用例
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{4, 5, 6, 7, 8}
	expected := []int{1, 2, 3}
	result := slice.DiffSlice(slice1, slice2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// 反向测试用例
	slice1 = []int{1, 2, 3, 4, 5}
	slice2 = []int{1, 2, 3, 4, 5}
	expected = []int{}
	result = slice.DiffSlice(slice1, slice2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// 测试别的类型
	slice3 := []string{"1", "2", "3", "4", "5"}
	slice4 := []string{"4", "5", "6", "7", "8"}
	expected2 := []string{"1", "2", "3"}
	result2 := slice.DiffSlice(slice3, slice4)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Expected %v, but got %v", expected2, result2)
	}

	// 测试interface{}类型
	slice5 := []interface{}{"1", "2", "3", "4", "5"}
	slice6 := []interface{}{"4", "5", "6", "7", "8"}
	expected3 := []interface{}{"1", "2", "3"}
	result3 := slice.DiffSlice(slice5, slice6)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Expected %v, but got %v", expected3, result3)
	}
}
