// 添加根路径
package main

import (
	"net/http"
	"os"
	"runtime"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	const (
		baseURL      = "http://192.168.31.253:58050/api/v1"
		apiPath      = "/rootPath/add"
		duration     = 10 * time.Second
		concurrency  = 100
		totalRPS     = 1000
	)

	headers := http.Header{
		"nas-user":     []string{"ca5a80f4-1ec5-42f4-8eec-4c07e5bbcd6f"},
		"Content-Type": []string{"application/json"},
	}

	payload := `{
		"partitionLabel": "666",
		"partitionUUID": "985ab745-0818-421d-b3d5-c44587c16f33",
		"path": "/some",
		"storageType":1
	}` // 表格数据

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
	for res := range attacker.Attack(targeter, rate, duration, "AddRootPath") {
		metrics.Add(res)
		// fmt.Printf("Response Body:\n%s\n\n", string(res.Body))

	}
	metrics.Close()

	vegeta.NewTextReporter(&metrics).Report(os.Stdout)
	transport.CloseIdleConnections()
}
