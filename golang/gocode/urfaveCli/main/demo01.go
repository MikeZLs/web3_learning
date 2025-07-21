package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	//test01()
	//test02()
	calculatorDemo()
}

func test01() {
	// 1. 创建一个新的 cli.Command 作为根命令 (v3 中 App 实际上也是一个 Command)
	cmd := &cli.Command{
		Name:  "greet",
		Usage: "跟世界打个招呼",
		// 2. 定义 Action，注意新的函数签名！
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// ctx 是标准库的 context.Context，用于处理超时、取消等。
			// cmd 是当前执行的命令对象，用于获取标志、参数等。
			fmt.Println("hello world!")
			return nil
		},
	}

	// 3. 运行根命令
	//    注意：Run 方法现在也需要一个 context.Context。
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func test02() {
	cmd := &cli.Command{
		Name:  "greet",
		Usage: "跟世界打个招呼",
		// 1. 定义标志
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Value:   "世界", // 默认值
				Usage:   "指定问候的对象",
				Aliases: []string{"n"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// 2. 从 cmd 对象获取标志的值
			//    注意：v3 中 Value() 方法返回 any (interface{})，需要进行类型断言。
			name := cmd.Value("name").(string)
			fmt.Printf("你好, %s!\n", name)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
