package testDemo

/*
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串s，判断字符串是否有效

有效字符串需满足：
左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。
*/

func Test02(s string) bool {
	// 创建一个栈，存放左括号
	stack := []rune{}

	// 定义括号的映射关系
	mapping := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, ch := range s {
		// 如果是右括号
		if open, exists := mapping[ch]; exists {
			// 如果栈为空 或 栈顶不匹配
			if len(stack) == 0 || stack[len(stack)-1] != open {
				return false
			}
			// 弹出栈顶
			stack = stack[:len(stack)-1]
		} else {
			// 是左括号，压入栈中
			stack = append(stack, ch)
		}
	}

	// 栈为空说明全部匹配成功
	return len(stack) == 0

}
