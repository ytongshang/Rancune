package util

import (
	"fmt"
)

func InSlice(slice []interface{}, value interface{}) bool {
	for _, vv := range slice {
		if vv == value {
			return true
		}
	}
	return false
}

func SliceInsert(slice []interface{}, index int, value interface{}) []interface{} {
	CheckRange(len(slice)+1, index)
	tmp := append([]interface{}{}, slice[index:]...)
	return append(append(slice[0:index], value), tmp...)
}

func SliceRemove(slice []interface{}, index int) []interface{} {
	CheckRange(len(slice), index)
	return append(slice[:index], slice[index+1:]...)
}

func SliceUpdate(slice []interface{}, index int, value interface{}) []interface{} {
	CheckRange(len(slice), index)
	slice[index] = value
	return slice
}

func SliceFilter(slice []interface{}, filter func(t interface{}) bool) (dslice []interface{}) {
	for _, v := range slice {
		dslice = append(dslice, filter(v))
	}
	return
}

func SliceMap(slice []interface{}, mapper func(t interface{}) interface{}) (ftslice []interface{}) {
	for _, vv := range slice {
		ftslice = append(ftslice, mapper(vv))
	}
	return
}

func SliceUnique(slice []interface{}) (uniqueslice []interface{}) {
	for _, v := range slice {
		if !InSlice(uniqueslice, v) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}

func CheckRange(size int, index int) int {
	if index < 0 && index >= size {
		panic(fmt.Sprintf("size = %d, but index = %d", size, index))
	}
	return index
}
