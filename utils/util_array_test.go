package utils

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	// 示例：将整数切片映射到它们的平方
	ints := []int{1, 2, 3, 4, 5}
	squared := Map(ints, func(x int) int {
		return x * x
	})
	fmt.Println(squared) // 输出: [1 4 9 16 25]

	// 示例：将字符串切片映射到它们的长度
	strs := []string{"hello", "world", "golang"}
	lengths := Map(strs, func(s string) int {
		return len(s)
	})
	fmt.Println(lengths) // 输出: [5 5 6]
}

func TestFilter(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 定义一个过滤条件，例子：过滤出偶数
	isEven := func(n int) bool {
		return n%2 == 0
	}

	// 使用filter函数
	filteredArr := Filter(arr, isEven)

	fmt.Println(filteredArr) // [2 4 6 8 10]
}
