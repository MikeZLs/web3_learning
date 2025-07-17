package main

import "web3_learning/golang/test/testDemo"

func main() {
	//testDemo.Test01()

	// // Test02
	// inputs := []string{
	// 	"()",
	// 	"()[]{}",
	// 	"(]",
	// 	"([)]",
	// 	"{[]}",
	// 	"",
	// 	"{[(])}",
	// }

	// for _, s := range inputs {
	// 	fmt.Printf("输入: %-10s  是否有效: %v\n", s, testDemo.Test02(s))
	// }

	// // Test03
	// tests := [][]string{
	// 	{"flower", "flow", "flight"},
	// 	{"dog", "racecar", "car"},
	// 	{"", "abc"},
	// 	{"abc"},
	// 	{"interview", "intermediate", "internal"},
	// }

	// for _, strs := range tests {
	// 	fmt.Printf("输入: %v\n输出: %q\n\n", strs, testDemo.Test03(strs))
	// }

	// // Test04
	// // 示例数组，必须是升序数组（题目要求）
	// nums := []int{1, 1, 2, 2, 3, 3, 3, 4}

	// fmt.Println(nums)

	// // 调用函数
	// k := testDemo.Test04(nums)
	// fmt.Println(nums)

	// // 打印结果（只打印前 k 个去重后的元素）
	// fmt.Printf("去重后长度：%d\n", k)
	// fmt.Printf("去重后的数组：%v\n", nums[:k])

	// // Test05
	// num := 5
	// fmt.Println("原始值：", num)
	// testDemo.Test05(&num) // 通过指针传递 num 的地址
	// fmt.Println("修改后的值：", num)

	// // Test06
	// nums := []int{1, 2, 3, 4, 5}
	// fmt.Println("原始数组：", nums)
	// testDemo.Test06(&nums)
	// fmt.Println("修改后的数组：", nums)

	// // Test07
	// testDemo.Test07()

	// // Test08
	// tasks := []testDemo.Task{
	// 	func() {
	// 		fmt.Println("任务1")
	// 		time.Sleep(time.Second * 2)
	// 	},
	// 	func() {
	// 		fmt.Println("任务2")
	// 		time.Sleep(time.Second)
	// 	},
	// 	func() {
	// 		fmt.Println("任务3")
	// 		time.Sleep(time.Second * 3)
	// 	},
	// }

	// testDemo.Test08(tasks)

	// // Test09
	// testDemo.Test09()

	// // Test10
	// testDemo.Test10()

	// // Test11
	// testDemo.Test11()

	// // Test12
	// testDemo.Test12()

	// // Test13
	// testDemo.Test13()

	// // Test14
	//testDemo.Test14()

	//	Test15
	//arr := []int{1, 9, 0, 7}
	//newArr := testDemo.Test15(arr)
	//fmt.Println(newArr)

	//// Test16
	//testDemo.Test16()

	// Test17
	testDemo.Test17()
}
