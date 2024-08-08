package utils

type MapperFunc[T any, U any] func(T) U

// Map 函数对输入切片应用 mapper 函数，并返回一个结果切片
func Map[T any, U any](input []T, mapper MapperFunc[T, U]) []U {
	output := make([]U, len(input))
	for i, v := range input {
		output[i] = mapper(v)
	}
	return output
}

// 定义一个函数类型，用于过滤条件
type Predicate func(int) bool

// 实现filter功能
func Filter(arr []int, predicate Predicate) []int {
	var result []int
	for _, v := range arr {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
