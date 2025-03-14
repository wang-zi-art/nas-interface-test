package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	baseURL  = "http://192.168.31.253:50080"
	nasToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjE1MDcwNDU3MDQ2IiwicGFzc3dvcmRWZXJzaW9uIjoiYzJiOTEyODAtOTE2OS00ZGQ4LWIyMjctMWE2ZjgxNzY1NmM5IiwicmFuZG9tS2V5IjoiMmI0Mzc1YzItMDI4OS00YzBmLWJjZTMtMjExYWIzYzNjODg5IiwiaXNzIjoiTGluY1VzZXJNYW5hZ2VyIn0.gh4-NcQm39I4D6URDFLuxMhAn06ELqG_JTw6uH_B7Uo"
)

func main() {
	fmt.Println("===== 开始执行 NAS 压测脚本 =====")

	files, err := filepath.Glob("*.go")
	if err != nil {
		fmt.Println("错误: 无法读取当前目录文件", err)
		return
	}

	for _, file := range files {
		if file == "main.go" {
			continue // 跳过自身
		}

		fmt.Printf("\n===== 执行: %s =====\n", file)

		cmd := exec.Command("go", "run", file)
		cmd.Env = append(os.Environ(),
			"BASE_URL="+baseURL,
			"NAS_TOKEN="+nasToken,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		start := time.Now()
		err := cmd.Run()
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("错误: 执行 %s 失败: %v\n", file, err)
		} else {
			fmt.Printf("完成: %s (耗时: %s)\n", file, duration)
		}
	}

	fmt.Println("===== 所有压测脚本执行完毕 =====")
}
