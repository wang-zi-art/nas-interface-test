package main

import (
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	// ========== 系统级优化配置 ==========
	runtime.GOMAXPROCS(runtime.NumCPU())

	// ========== 压测参数配置 ==========
	const (
		baseURL      = "http://192.168.31.253:58050/api/v1"
		apiPath      = "/AI/face/facesByAlias" // 来自xlsx的接口路径
		duration     = 10 * time.Second
		concurrency  = 100   // 并发数
		totalRPS     = 1000  // 总请求率 = 100并发 * 10次/秒
	)

	// ========== 请求参数配置 ==========
	headers := http.Header{
		"nas-user":     []string{"ca5a80f4-1ec5-42f4-8eec-4c07e5bbcd6f"},
		"Content-Type": []string{"application/json"},
		"Connection":   []string{"keep-alive"}, // 连接复用
	}

	// 来自xlsx的Body参数（原样使用）
	payload := `{"alias": "ccd"}`

	// ========== 高性能HTTP客户端配置 ==========
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,  // 保持连接时间
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          1000,      // 最大空闲连接数
		MaxIdleConnsPerHost:   500,       // 每个目标主机的最大空闲连接
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,  // 单次请求超时
		}),
		vegeta.Workers(uint64(concurrency)),
		vegeta.KeepAlive(true),
	)

	// ========== 构建完整请求目标 ==========
	fullURL := baseURL + apiPath
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    fullURL,
		Body:   []byte(payload),
		Header: headers,
	})

	// ========== 执行压测并收集指标 ==========
	rate := vegeta.Rate{Freq: totalRPS, Per: time.Second}
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "FacesByAlias") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))

	}
	metrics.Close()

	// ========== 生成报告并清理资源 ==========
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)
	transport.CloseIdleConnections()
}