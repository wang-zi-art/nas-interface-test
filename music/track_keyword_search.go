package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	baseURL      = "http://192.168.31.253:50080"
	apiPath      = "/service/audioPlayer/api/v1/secure/track/find/keyword"
	testDuration = 10 * time.Second
	workers      = 100
	totalRPS     = 1000
	nasToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjE1MDcwNDU3MDQ2IiwicGFzc3dvcmRWZXJzaW9uIjoiYzJiOTEyODAtOTE2OS00ZGQ4LWIyMjctMWE2ZjgxNzY1NmM5IiwicmFuZG9tS2V5IjoiMmI0Mzc1YzItMDI4OS00YzBmLWJjZTMtMjExYWIzYzNjODg5IiwiaXNzIjoiTGluY1VzZXJNYW5hZ2VyIn0.gh4-NcQm39I4D6URDFLuxMhAn06ELqG_JTw6uH_B7Uo"
)

var (
	requestBody = `{
		"findType": "track",
		"keyword": "借口"
	}`
)

func main() {
	// 配置HTTP传输参数
	transport := &http.Transport{
		MaxIdleConnsPerHost: 500,
	}

	// 初始化攻击器
	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}),
		vegeta.Workers(uint64(workers)),
		vegeta.KeepAlive(true),
	)

	// 配置固定请求目标
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    baseURL + apiPath,
		Body:   []byte(requestBody),
		Header: http.Header{
			"nas-token":    []string{nasToken},
			"Content-Type": []string{"application/json"},
		},
	})

	// 执行压测
	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: totalRPS, Per: time.Second}

	for res := range attacker.Attack(targeter, rate, testDuration, "关键词搜索压测") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))
	}
	metrics.Close()

	// 输出统计报告
	fmt.Println("\n================ 压测摘要 ================")
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)

	// 输出关键指标
	fmt.Println("\n======= 核心指标 =======")
	fmt.Printf("成功率: %.2f%%\n", metrics.Success*100)
	fmt.Printf("总请求数: %d\n", metrics.Requests)
	fmt.Printf("平均延迟: %.2fms\n", metrics.Latencies.Mean.Seconds()*1000)

	// 清理资源
	transport.CloseIdleConnections()
}
