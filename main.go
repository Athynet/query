package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 命令行参数解析
	var (
		inputFile   = flag.String("i", "test.csv", "输入CSV文件路径")
		outputFile  = flag.String("o", "output.csv", "输出CSV文件路径")
		keyFile     = flag.String("k", "private.pem", "RSA私钥文件路径")
		concurrency = flag.Int("c", 4, "并发处理的goroutine数量")
	)
	flag.Parse()

	// 处理输入文件 - 如果test.csv不存在，尝试使用text.csv
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		if *inputFile == "test.csv" {
			*inputFile = "text.csv"
			fmt.Printf("使用text.csv作为输入文件\n")
		} else {
			fmt.Fprintf(os.Stderr, "错误: 输入文件 %s 不存在\n", *inputFile)
			os.Exit(1)
		}
	}

	// 加载RSA私钥
	privateKey, err := LoadPrivateKey(*keyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载私钥失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("已加载私钥: %s\n", *keyFile)

	// 直接使用字符串变量定义模板，使用fmt.Sprintf支持的%s占位符
	template := "trade_no=%s&version=1.0"
	fmt.Printf("使用模板: %s\n", template)

	// 处理CSV记录
	fmt.Println("处理CSV记录，进行RSA-PSS签名...")
	signFunc := func(data []byte) (string, error) {
		return RSA_PSS_Sign(privateKey, data)
	}

	// 使用流式并发处理
	fmt.Printf("使用%d个并发goroutine处理...\n", *concurrency)
	if err := ProcessCSVStream(*inputFile, *outputFile, signFunc, template, *concurrency); err != nil {
		fmt.Fprintf(os.Stderr, "处理CSV文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("处理完成！")
}
