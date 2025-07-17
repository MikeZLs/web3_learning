package testDemo

// 题目：给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一

func Test15(arr []int) []int {
	//for i := len(arr) - 1; i >= 0; i-- {
	//	arr[i]++
	//	if arr[i] < 10 {
	//		return arr
	//	}
	//	arr[i] = 0
	//}
	//
	//result := make([]int, len(arr)+1)
	//result[0] = 1
	//
	//return result

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] < 9 {
			arr[i]++
			return arr
		}
		arr[i] = 0
	}
	return append([]int{1}, arr...)

}
