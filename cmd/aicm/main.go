package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cocotyty/aicm/internal/config"
	"github.com/cocotyty/aicm/internal/generator"
	"github.com/cocotyty/aicm/internal/git"
)

func main() {
	// 初始化日志
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 处理配置命令
	if len(os.Args) > 1 && os.Args[1] == "config" {
		if len(os.Args) < 4 {
			fmt.Println("Usage: aicm config set KEY=VALUE")
			os.Exit(1)
		}

		if os.Args[2] != "set" {
			fmt.Println("Unknown config command:", os.Args[2])
			os.Exit(1)
		}

		parts := strings.SplitN(os.Args[3], "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid config format. Use KEY=VALUE")
			os.Exit(1)
		}

		key := strings.ToUpper(parts[0])
		value := parts[1]

		if err := config.SetConfig(key, value); err != nil {
			fmt.Println("Failed to set config:", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set %s\n", key)
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		fmt.Println("Please configure first using:")
		fmt.Println("  aicm config set LLM_API_KEY=your_key")
		fmt.Println("  aicm config set LLM_MODEL=your_model")
		fmt.Println("  aicm config set LLM_API_URL=your_url")
		os.Exit(1)
	}

	// 获取git变更
	changes, err := git.GetChanges()
	if err != nil {
		fmt.Println("Failed to get git changes:", err)
		os.Exit(1)
	}

	// 生成commit message
	msg, err := generator.GenerateCommitMessage(cfg, changes)
	if err != nil {
		fmt.Println("Failed to generate commit message:", err)
		os.Exit(1)
	}

	fmt.Println("Generated commit message:")
	fmt.Println(msg)
}
