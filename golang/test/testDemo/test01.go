package testDemo

import "fmt"

/*
	给定一个非空整数数组，除了某个元素只出现一次以外，
	其余每个元素均出现两次。
	找出那个只出现了一次的元素。
	可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，
	例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素
*/

func Test01() {
	// 创建一个非空整数数组
	arr := []int{1, 2, 4, 1, 2, 3, 4}

	// 创建一个 map 用于记录每个元素出现的次数
	countMap := make(map[int]int)

	// 遍历数组
	for _, num := range arr {
		countMap[num]++
	}

	fmt.Println(countMap)

	// 遍历countMap
	for k, v := range countMap {
		if v == 1 {
			fmt.Println("只出现一次的元素为：", k)
		}
	}

}
