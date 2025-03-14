// 刷新歌单压测
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
	apiPath      = "/service/audioPlayer/api/v1/secure/playlist/refresh"
	testDuration = 10 * time.Second
	workers      = 100
	totalRPS     = 1000
	nasToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjE1MDcwNDU3MDQ2IiwicGFzc3dvcmRWZXJzaW9uIjoiYzJiOTEyODAtOTE2OS00ZGQ4LWIyMjctMWE2ZjgxNzY1NmM5IiwicmFuZG9tS2V5IjoiMmI0Mzc1YzItMDI4OS00YzBmLWJjZTMtMjExYWIzYzNjODg5IiwiaXNzIjoiTGluY1VzZXJNYW5hZ2VyIn0.gh4-NcQm39I4D6URDFLuxMhAn06ELqG_JTw6uH_B7Uo"
)

var requestBody = `{"playlistId": 8}`

func main() {
	transport := &http.Transport{MaxIdleConnsPerHost: 500}

	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}),
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

	for res := range attacker.Attack(targeter, rate, testDuration, "刷新歌单压测") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))

	}
	metrics.Close()

	fmt.Println("\n======= 关键指标 =======")
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)

	fmt.Printf("平均延迟: %.2fms\n", metrics.Latencies.Mean.Seconds()*1000)
	fmt.Printf("最大QPS: %.1f\n", metrics.Rate)
	transport.CloseIdleConnections()
}
