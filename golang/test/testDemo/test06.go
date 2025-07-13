package testDemo

/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/

// 传入一个整数切片的指针
func Test06(nums *[]int) {
	// 遍历切片
	for i := range *nums {
		// 将每个元素乘以2
		(*nums)[i] *= 2
	}
}
