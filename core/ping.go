package core

import (
	"flag"
	"fmt"
	"log"
	"net"
	"speedgo/commands"
	"strings"
	"sync"
	"time"
)

type PingResult struct {
	Target string
	RTTs   []time.Duration
	MinRTT time.Duration
	MaxRTT time.Duration
	AvgRTT time.Duration
	Lost   int
}

func RunPing(args []string) {
	if err := commands.PingCmd.Parse(args); err != nil {
		fmt.Println("Error parsing ping arguments:", err)
		commands.PingCmd.PrintDefaults()
		return
	}

	targets := commands.PingCmd.Lookup("targets").Value.String()
	count := commands.PingCmd.Lookup("count").Value.(flag.Getter).Get().(int)
	timeout := commands.PingCmd.Lookup("timeout").Value.(flag.Getter).Get().(time.Duration)
	concurrency := commands.PingCmd.Lookup("concurrency").Value.(flag.Getter).Get().(int)
	verbose := commands.PingCmd.Lookup("verbose").Value.(flag.Getter).Get().(bool)

	targetList := splitTargets(targets)
	if len(targetList) == 0 {
		fmt.Println("Error: No valid targets provided")
		commands.PingCmd.PrintDefaults()
		return
	}

	fmt.Printf("Starting ping test to %d targets...\n", len(targetList))
	results := pingTargets(targetList, count, timeout, concurrency, verbose)

	// 输出结果
	for _, result := range results {
		fmt.Printf("\nPing results for %s:\n", result.Target)
		if len(result.RTTs) == 0 {
			fmt.Printf("  No responses (all packets lost)\n")
		} else {
			fmt.Printf("  Min RTT: %v\n", result.MinRTT)
			fmt.Printf("  Max RTT: %v\n", result.MaxRTT)
			fmt.Printf("  Avg RTT: %v\n", result.AvgRTT)
			fmt.Printf("  Packet Loss: %d/%d\n", result.Lost, count)
		}
	}
}

func splitTargets(targets string) []string {
	var result []string
	for _, target := range splitAndTrim(targets, ",") {
		if net.ParseIP(target) != nil || len(target) > 0 {
			result = append(result, target)
		}
	}
	return result
}

func splitAndTrim(input, sep string) []string {
	parts := strings.Split(input, sep)
	var result []string
	for _, item := range parts {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// 并发 Ping 实现
func pingTargets(targets []string, count int, timeout time.Duration, concurrency int, verbose bool) []PingResult {
	var wg sync.WaitGroup
	results := make([]PingResult, len(targets))
	targetChan := make(chan int, len(targets))
	for i := range targets {
		targetChan <- i
	}
	close(targetChan)

	//sem := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range targetChan {
				results[idx] = pingTarget(targets[idx], count, timeout, verbose)
				if verbose {
					fmt.Printf("Completed ping to %s\n", targets[idx])
				}
			}
		}()
	}

	wg.Wait()
	return results
}

func pingTarget(target string, count int, timeout time.Duration, verbose bool) PingResult {
	result := PingResult{Target: target, RTTs: []time.Duration{}, Lost: 0}
	var totalRTT time.Duration

	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("ip4:icmp", target, timeout)
		if err != nil {
			if verbose {
				fmt.Printf("Ping %s failed: %v\n", target, err)
			}
			result.Lost++
			continue
		}

		elapsed := time.Since(start)
		result.RTTs = append(result.RTTs, elapsed)
		totalRTT += elapsed
		err = conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}

		if verbose {
			fmt.Printf("Ping %s: RTT = %v\n", target, elapsed)
		}

		time.Sleep(time.Millisecond * 500) // 间隔
	}

	// 计算统计值
	if len(result.RTTs) > 0 {
		result.MinRTT = minRTT(result.RTTs)
		result.MaxRTT = maxRTT(result.RTTs)
		result.AvgRTT = totalRTT / time.Duration(len(result.RTTs))
	}
	return result
}

func minRTT(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	minVal := durations[0]
	for _, d := range durations {
		if d < minVal {
			minVal = d
		}
	}
	return minVal
}

func maxRTT(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	maxVal := durations[0]
	for _, d := range durations {
		if d > maxVal {
			maxVal = d
		}
	}
	return maxVal
}
