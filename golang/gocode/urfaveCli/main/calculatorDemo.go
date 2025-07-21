package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"strconv"
)

func calculatorDemo() {
	cmd := &cli.Command{
		Name:  "calc",
		Usage: "一个简单的计算器",
		// 1. 定义子命令列表
		Commands: []*cli.Command{
			{
				Name:    "add",
				Usage:   "计算两个数字的和",
				Aliases: []string{"a"},
				// 2. 子命令的 Action
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// 3. 从 cmd.Args() 获取参数
					arg1, _ := strconv.Atoi(cmd.Args().Get(0))
					arg2, _ := strconv.Atoi(cmd.Args().Get(1))
					fmt.Printf("%d + %d = %d\n", arg1, arg2, arg1+arg2)
					return nil
				},
			},
			{
				Name:    "subtract",
				Usage:   "计算两个数字的差",
				Aliases: []string{"s"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					arg1, _ := strconv.Atoi(cmd.Args().Get(0))
					arg2, _ := strconv.Atoi(cmd.Args().Get(1))
					fmt.Printf("%d - %d = %d\n", arg1, arg2, arg1-arg2)
					return nil
				},
			},
			{
				Name:    "product",
				Usage:   "计算两个数字的乘积",
				Aliases: []string{"p"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					arg1, _ := strconv.Atoi(cmd.Args().Get(0))
					arg2, _ := strconv.Atoi(cmd.Args().Get(1))
					fmt.Printf("%d * %d = %d\n", arg1, arg2, arg1*arg2)
					return nil
				},
			},
			{
				Name:    "quotient",
				Usage:   "计算两个数字相除的商",
				Aliases: []string{"q"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					arg1, _ := strconv.Atoi(cmd.Args().Get(0))
					arg2, _ := strconv.Atoi(cmd.Args().Get(1))
					fmt.Printf("%d / %d = %d\n", arg1, arg2, arg1/arg2)
					return nil
				},
			},
			{
				Name:    "modulo",
				Usage:   "计算两个数字取余",
				Aliases: []string{"m"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					arg1, _ := strconv.Atoi(cmd.Args().Get(0))
					arg2, _ := strconv.Atoi(cmd.Args().Get(1))
					fmt.Printf("%d % %d = %d\n", arg1, arg2, arg1%arg2)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
