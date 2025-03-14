// 专辑查询压测
package main

import (
	"fmt"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	baseURL      = "http://192.168.31.253:50080"
	apiPath      = "/service/audioPlayer/api/v1/secure/album/find/albumId"
	testDuration = 10 * time.Second
	workers      = 100
	totalRPS     = 1000
	nasToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjE1MDcwNDU3MDQ2IiwicGFzc3dvcmRWZXJzaW9uIjoiYzJiOTEyODAtOTE2OS00ZGQ4LWIyMjctMWE2ZjgxNzY1NmM5IiwicmFuZG9tS2V5IjoiMmI0Mzc1YzItMDI4OS00YzBmLWJjZTMtMjExYWIzYzNjODg5IiwiaXNzIjoiTGluY1VzZXJNYW5hZ2VyIn0.gh4-NcQm39I4D6URDFLuxMhAn06ELqG_JTw6uH_B7Uo"
)

var (
	requestBody = `{
		"albumId": 5,
		"orderParam": {
			"ascending": true,
			"field": "alphabetical"
		},
		"page": {
			"pageIndex": 1,
			"pageSize": 20
		}
	}`
)

func main() {
	transport := &http.Transport{MaxIdleConnsPerHost: 500}

	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{Transport: transport, Timeout: 30*time.Second}),
		vegeta.Workers(uint64(workers)),
	)

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    baseURL + apiPath,
		Body:   []byte(requestBody),
		Header: http.Header{
			"nas-token":    []string{nasToken},
			"Content-Type": []string{"application/json"},
		},
	})

	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: totalRPS, Per: time.Second}

	for res := range attacker.Attack(targeter, rate, testDuration, "专辑查询压测") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))
	}
	metrics.Close()

	fmt.Println("\n======= 性能指标 =======")
	fmt.Printf("请求成功率: %.2f%%\n", metrics.Success*100)
	fmt.Printf("平均延迟: %.2fms\n", metrics.Latencies.Mean.Seconds()*1000)
	transport.CloseIdleConnections()
}
