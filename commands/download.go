package commands

import (
	"flag"
	"time"
)

var DownloadCmd = flag.NewFlagSet("download", flag.ExitOnError)

func init() {
	DownloadCmd.String("url", "", "URL to download from (required)")
	DownloadCmd.Duration("duration", time.Second*30, "Maximum download duration")
	DownloadCmd.Int("concurrency", 4, "Number of concurrent download chunks")
	DownloadCmd.String("output", "", "Output file path (optional)")
	DownloadCmd.Bool("verbose", false, "Enable detailed output")
}
