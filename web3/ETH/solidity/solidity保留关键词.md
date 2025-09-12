以下是Solidity中的主要关键词，按类别整理：

## **数据类型关键词**
- `bool` - 布尔类型
- `int` / `uint` - 有符号/无符号整数
- `int8` ~ `int256` / `uint8` ~ `uint256` - 指定位数的整数
- `fixed` / `ufixed` - 定点数（实验性）
- `address` - 地址类型
- `bytes` - 动态字节数组
- `bytes1` ~ `bytes32` - 固定长度字节数组
- `string` - 字符串类型

## **存储位置关键词**
- `storage` - 状态变量存储
- `memory` - 函数内存储
- `calldata` - 函数参数存储（只读）

## **可见性修饰符**
- `public` - 公开可见
- `private` - 仅当前合约可见
- `internal` - 当前合约及继承合约可见
- `external` - 仅外部调用可见

## **状态可变性修饰符**
- `pure` - 不读取也不修改状态
- `view` - 只读取状态，不修改
- `payable` - 可接收以太币
- `constant` - 常量（已废弃，用view替代）

## **继承和多态关键词**
- `virtual` - 函数可被重写
- `override` - 重写父合约函数
- `abstract` - 抽象合约
- `interface` - 接口声明

## **合约结构关键词**
- `contract` - 合约声明
- `library` - 库声明
- `interface` - 接口声明
- `struct` - 结构体
- `enum` - 枚举类型
- `event` - 事件声明
- `error` - 自定义错误
- `modifier` - 修饰器

## **函数相关关键词**
- `function` - 函数声明
- `constructor` - 构造函数
- `fallback` - 回退函数
- `receive` - 接收以太币函数
- `returns` - 返回值声明
- `return` - 返回语句

## **继承关键词**
- `is` - 继承声明
- `super` - 调用父合约

## **控制流关键词**
- `if` / `else` - 条件语句
- `while` - while循环
- `for` - for循环
- `do` - do-while循环
- `break` - 跳出循环
- `continue` - 继续下次循环

## **异常处理关键词**
- `require` - 条件检查
- `assert` - 断言检查
- `revert` - 主动回退
- `try` / `catch` - 异常捕获

## **其他重要关键词**
- `pragma` - 编译器指令
- `import` - 导入文件
- `using` - 库函数绑定
- `assembly` - 内联汇编
- `emit` - 触发事件
- `delete` - 删除变量
- `new` - 创建合约实例
- `this` - 当前合约引用
- `selfdestruct` - 销毁合约

## **保留字**
- `after`, `alias`, `apply`, `auto`, `case`, `copyof`, `default`, `define`, `final`, `immutable`, `implements`, `in`, `inline`, `let`, `macro`, `match`, `mutable`, `null`, `of`, `partial`, `promise`, `reference`, `relocatable`, `sealed`, `sizeof`, `static`, `supports`, `switch`, `typedef`, `typeof`, `var`
