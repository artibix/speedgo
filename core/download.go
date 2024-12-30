// Package core download.go
package core

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"speedgo/commands"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// DownloadConfig stores download test configuration
type DownloadConfig struct {
	Duration    time.Duration
	Concurrency int
	Verbose     bool
}

// DownloadStats stores download speed statistics
type DownloadStats struct {
	BytesReceived int64
	Duration      time.Duration
	Speed         float64 // Speed in Mbps
	Error         error
}

// Default test files from various CDNs
var defaultTestFiles = []string{
	"https://speed.cloudflare.com/__down?bytes=25000000", // 25MB test file
	"https://cdn.jsdelivr.net/gh/librespeed/speedtest-files@master/random4000x4000.jpg",
	"https://proof.ovh.net/files/100Mb.dat",
}

func RunDownload(ctx context.Context, args []string) error {
	config, err := parseDownloadConfig(args)
	if err != nil {
		return fmt.Errorf("parsing download config: %w", err)
	}

	fmt.Printf("Starting download speed test (Duration: %v, Concurrent streams: %d)\n",
		config.Duration, config.Concurrency)

	stats := measureDownloadSpeed(ctx, config)
	printDownloadResults(stats)

	return nil
}

func measureDownloadSpeed(ctx context.Context, config *DownloadConfig) DownloadStats {
	var totalBytes int64
	start := time.Now()

	// Create channels for coordination
	errChan := make(chan error, config.Concurrency)
	bytesChan := make(chan int64, config.Concurrency)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, config.Duration)
	defer cancel()

	// Start concurrent downloads
	var wg sync.WaitGroup
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			downloadWorker(ctx, workerID, config, bytesChan, errChan)
		}(i)
	}

	// Start progress monitoring in separate goroutine
	go func() {
		if config.Verbose {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					current := atomic.LoadInt64(&totalBytes)
					duration := time.Since(start)
					speed := float64(current*8) / (1000 * 1000 * duration.Seconds()) // Mbps
					fmt.Printf("\rCurrent speed: %.2f Mbps", speed)
				}
			}
		}
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(bytesChan)
		close(errChan)
	}()

	// Process results
	var lastError error
	for {
		select {
		case bytes, ok := <-bytesChan:
			if !ok {
				duration := time.Since(start)
				return DownloadStats{
					BytesReceived: totalBytes,
					Duration:      duration,
					Speed:         float64(totalBytes*8) / (1000 * 1000 * duration.Seconds()),
					Error:         lastError,
				}
			}
			atomic.AddInt64(&totalBytes, bytes)

		case err := <-errChan:
			if err != nil {
				lastError = err
			}
		}
	}
}

func downloadWorker(ctx context.Context, id int, config *DownloadConfig,
	bytesChan chan<- int64, errChan chan<- error) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Choose a random test file
			url := defaultTestFiles[id%len(defaultTestFiles)]

			if err := downloadChunk(ctx, url, bytesChan); err != nil {
				errChan <- fmt.Errorf("worker %d error: %w", id, err)
				time.Sleep(time.Second) // Back off on error
				continue
			}
		}
	}
}

func downloadChunk(ctx context.Context, url string, bytesChan chan<- int64) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			bytesChan <- int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading response: %w", err)
		}
	}

	return nil
}

func parseDownloadConfig(args []string) (*DownloadConfig, error) {
	cmd := commands.DownloadCmd
	if err := cmd.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing arguments: %w", err)
	}

	return &DownloadConfig{
		Duration:    cmd.Lookup("duration").Value.(flag.Getter).Get().(time.Duration),
		Concurrency: cmd.Lookup("concurrency").Value.(flag.Getter).Get().(int),
		Verbose:     cmd.Lookup("verbose").Value.(flag.Getter).Get().(bool),
	}, nil
}

func printDownloadResults(stats DownloadStats) {
	fmt.Printf("\n\nDOWNLOAD TEST RESULTS\n")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total data received: %.2f MB\n", float64(stats.BytesReceived)/(1024*1024))
	fmt.Printf("Test duration: %.1f seconds\n", stats.Duration.Seconds())
	fmt.Printf("Average speed: %.2f Mbps\n", stats.Speed)
	if stats.Error != nil {
		fmt.Printf("Errors encountered: %v\n", stats.Error)
	}
	fmt.Println(strings.Repeat("=", 50))
}
