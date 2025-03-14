// 获取设备信息压测
package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	const (
		baseURL      = "https://tunnel-manage.techvision-nas.com:55000"
		apiPath      = "/api/v2/client/getDeviceInfo"
		duration     = 10 * time.Second
		concurrency  = 100
		totalRPS     = 1000
	)

	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}

	payload := `{
		"accessID": "JLWw06DKMPFG5F25ZvtN8"
	}`

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

	for res := range attacker.Attack(targeter, rate, duration, "GetDeviceInfo") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))
	}

	metrics.Close()

	// 输出 vegeta 统计信息
	fmt.Println("\n================ Vegeta Summary ================")
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)

	transport.CloseIdleConnections()
}
