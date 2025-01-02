package mobile

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type PingResult struct {
	Target  string  `json:"target"`
	MinRTT  float64 `json:"minRTT"`
	MaxRTT  float64 `json:"maxRTT"`
	AvgRTT  float64 `json:"avgRTT"`
	Lost    int     `json:"lost"`
	Success bool    `json:"success"`
	Error   string  `json:"error,omitempty"`
}

// tcpPing 使用TCP连接测试延迟
func tcpPing(target string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", target), timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return time.Since(start), nil
}

// Ping 封装给移动端调用的ping方法
func Ping(target string) (string, error) {
	timeout := time.Second * 3
	count := 3
	var rtts []time.Duration
	var lost int

	for i := 0; i < count; i++ {
		rtt, err := tcpPing(target, timeout)
		if err != nil {
			lost++
			continue
		}
		rtts = append(rtts, rtt)
		time.Sleep(time.Second)
	}

	result := PingResult{
		Target:  target,
		Success: len(rtts) > 0,
	}

	if len(rtts) > 0 {
		var minRTT, maxRTT, totalRTT time.Duration
		minRTT = rtts[0]
		maxRTT = rtts[0]

		for _, rtt := range rtts {
			totalRTT += rtt
			if rtt < minRTT {
				minRTT = rtt
			}
			if rtt > maxRTT {
				maxRTT = rtt
			}
		}

		// 转换成毫秒
		result.MinRTT = float64(minRTT.Microseconds()) / 1000
		result.MaxRTT = float64(maxRTT.Microseconds()) / 1000
		result.AvgRTT = float64(totalRTT.Microseconds()) / float64(len(rtts)) / 1000
	}

	result.Lost = lost

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
