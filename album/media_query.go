// 分页搜索媒体
package main

import (
	// "encoding/json"

	"net/http"
	"os"
	"runtime"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// 定义响应结构体
// type Response struct {
// 	Code    int    `json:"code"`
// 	Message string `json:"message"`
// }

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	const (
		baseURL      = "http://192.168.31.253:58050/api/v1"
		apiPath      = "/media/query"
		duration     = 10 * time.Second
		concurrency  = 100
		totalRPS     = 1000
	)

	headers := http.Header{
		"nas-user":     []string{"ca5a80f4-1ec5-42f4-8eec-4c07e5bbcd6f"},
		"Content-Type": []string{"application/json"},
	}

	payload := `{
"endTime": 99999999999999, "keyword": "", "page": 0, "pageSize": 80, "startTime": 0
	}` // 表格参数

	transport := &http.Transport{
		MaxIdleConnsPerHost: 500,
	}

	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{Transport: transport}),
		vegeta.Workers(uint64(concurrency)),
	)

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    baseURL + apiPath,
		Body:   []byte(payload),
		Header: headers,
	})

	rate := vegeta.Rate{Freq: totalRPS, Per: time.Second}
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "MediaQuery") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))
		// 解析 JSON 响应体
		// var response Response
		// if err := json.Unmarshal(res.Body, &response); err == nil {
		// 	fmt.Printf("Response: code=%d, message=%s\n", response.Code, response.Message)
		// } else {
		// 	fmt.Println("Failed to parse response JSON:", err)
		// }
	}
	metrics.Close()

	vegeta.NewTextReporter(&metrics).Report(os.Stdout)
	transport.CloseIdleConnections()
}
