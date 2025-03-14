// 文件来源检索压测
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
	apiPath      = "/service/audioPlayer/api/v1/secure/fileList/retrieve"
	testDuration = 10 * time.Second
	workers      = 100
	totalRPS     = 1000
	nasToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjE1MDcwNDU3MDQ2IiwicGFzc3dvcmRWZXJzaW9uIjoiYzJiOTEyODAtOTE2OS00ZGQ4LWIyMjctMWE2ZjgxNzY1NmM5IiwicmFuZG9tS2V5IjoiMmI0Mzc1YzItMDI4OS00YzBmLWJjZTMtMjExYWIzYzNjODg5IiwiaXNzIjoiTGluY1VzZXJNYW5hZ2VyIn0.gh4-NcQm39I4D6URDFLuxMhAn06ELqG_JTw6uH_B7Uo"
)

func main() {
	transport := &http.Transport{MaxIdleConnsPerHost: 500}

	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{
			Transport: transport,
			Timeout:   45 * time.Second, // 延长超时时间
		}),
		vegeta.Workers(uint64(workers)),
	)

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    baseURL + apiPath,
		Header: http.Header{
			"nas-token":    []string{nasToken},
			"Content-Type": []string{"application/json"},
		},
	})

	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: totalRPS, Per: time.Second}

	for res := range attacker.Attack(targeter, rate, testDuration, "文件来源检索压测") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))
	}
	metrics.Close()

	fmt.Println("\n======= 关键指标 =======")
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)

	fmt.Printf("平均吞吐量: %.2f req/s\n", metrics.Rate)
	fmt.Printf("95%%延迟: %.2fms\n", metrics.Latencies.P95.Seconds()*1000)
	transport.CloseIdleConnections()
}
