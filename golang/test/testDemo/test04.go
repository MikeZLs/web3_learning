package testDemo

// 给你一个非严格递增排列的数组 nums ，请你原地删除重复出现的元素，
// 使每个元素只出现一次 ，返回删除后数组的新长度。元素的相对顺序应该保持一致 。然后返回 nums 中唯一元素的个数。

// 在 Go 语言中，切片是引用类型，函数中对切片元素的修改会直接影响调用者看到的值，哪怕你只是返回了一个整数
func Test04(nums []int) int {
	k := 1

	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			nums[k] = nums[i]
			k++
		}
	}

	return k
}
