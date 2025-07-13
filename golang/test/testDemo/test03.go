package testDemo

/*
编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。

示例 1：
输入：strs = ["flower","flow","flight"]
输出："fl"

示例 2：
输入：strs = ["dog","racecar","car"]
输出：""
解释：输入不存在公共前缀。
*/

func Test03(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	// 以第一个字符串为基准
	prefix := strs[0]

	// 遍历其余的字符串
	for i := 1; i < len(strs); i++ {
		// 比较 prefix 和 strs[i] 的公共前缀
		for len(strs[i]) < len(prefix) || strs[i][:len(prefix)] != prefix {
			prefix = prefix[:len(prefix)-1] // 缩短 prefix
			if prefix == "" {
				return ""
			}
		}
	}

	return prefix
}
