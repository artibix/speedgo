package commands

import "flag"

var UploadCmd = flag.NewFlagSet("download", flag.ExitOnError)

func init() {
	UploadCmd.String("url", "http://speedtest.example.com", "Base URL of the speed test server (default: example server)")
	UploadCmd.Int("concurrency", 4, "Number of concurrent downloads (default: 4)")
	UploadCmd.Int("duration", 10, "Test duration in seconds")
	UploadCmd.Bool("verbose", false, "Enable detailed output")
}
