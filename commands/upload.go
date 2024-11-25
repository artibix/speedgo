package commands

import "flag"

var DownloadCmd = flag.NewFlagSet("download", flag.ExitOnError)

func init() {
	DownloadCmd.String("url", "http://speedtest.example.com", "Base URL of the speed test server (default: example server)")
	DownloadCmd.Int("concurrency", 4, "Number of concurrent downloads (default: 4)")
	DownloadCmd.Int("duration", 10, "Test duration in seconds")
	DownloadCmd.Bool("verbose", false, "Enable detailed output")
}
