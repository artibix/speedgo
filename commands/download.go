package commands

import "flag"

var UploadCmd = flag.NewFlagSet("upload", flag.ExitOnError)

func init() {
	UploadCmd.String("url", "", "URL to upload to (required)")
	UploadCmd.String("file", "", "Path to file for upload (optional, if not provided, generates test data)")
	UploadCmd.Int("size", 10, "Size of test data to generate in MB (ignored if file is specified)")
	UploadCmd.Int("chunk", 1, "Size of each chunk in MB (default: 1MB)")
	UploadCmd.Int("retries", 3, "Number of retries for failed chunks")
	UploadCmd.Bool("verbose", false, "Enable detailed output")
}
