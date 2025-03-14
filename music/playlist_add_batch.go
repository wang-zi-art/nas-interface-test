// 批量添加歌曲压测
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
	apiPath      = "/service/audioPlayer/api/v1/secure/playlist/add/batch"
	testDuration = 10 * time.Second
	workers      = 100
	totalRPS     = 1000
	nasToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjE1MDcwNDU3MDQ2IiwicGFzc3dvcmRWZXJzaW9uIjoiYzJiOTEyODAtOTE2OS00ZGQ4LWIyMjctMWE2ZjgxNzY1NmM5IiwicmFuZG9tS2V5IjoiMmI0Mzc1YzItMDI4OS00YzBmLWJjZTMtMjExYWIzYzNjODg5IiwiaXNzIjoiTGluY1VzZXJNYW5hZ2VyIn0.gh4-NcQm39I4D6URDFLuxMhAn06ELqG_JTw6uH_B7Uo"
)

var (
	requestBody = `{
		"playlistId": 2,
		"trackIdList": [2486,2491]
	}`
)

func main() {
	transport := &http.Transport{MaxIdleConnsPerHost: 500}

	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}),
		vegeta.Workers(uint64(workers)),
		vegeta.KeepAlive(true),
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

	for res := range attacker.Attack(targeter, rate, testDuration, "批量添加歌曲压测") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))
	}
	metrics.Close()

	fmt.Println("\n======= 核心指标 =======")
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)

	fmt.Printf("总成功率: %.2f%%\n", metrics.Success*100)
	fmt.Printf("请求总量: %d\n", metrics.Requests)
	transport.CloseIdleConnections()
}
