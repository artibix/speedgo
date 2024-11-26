package main

import (
	"fmt"
	"os"
	"speedgo/commands"
	"speedgo/core"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "ping", "p":
		pingCommand(args)
	case "download", "d":
		downloadCommand(args)
	case "upload", "u":
		uploadCommand(args)
	case "-h", "--help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage: speedgo <command> [options]")
	fmt.Println("\ncmd:")
	fmt.Println("  ping, p        Test network latency (ping multiple targets)")
	fmt.Println("  download, d    Test download speed")
	fmt.Println("  upload, u      Test upload speed")
	fmt.Println("\nExamples:")
	fmt.Println("  speedgo ping --targets=google.com --count=5")
	fmt.Println("  speedgo d --url=http://example.com/file.dat --duration=15")
	fmt.Println("  speedgo u --file=test.dat --url=http://example.com/upload")
	fmt.Println("\nHelp:")
	fmt.Println("  speedgo <command> -h    Show help for a specific command")
}

func pingCommand(args []string) {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		commands.PingCmd.Usage()
		os.Exit(0)
	}
	core.RunPing(args)
}

func downloadCommand(args []string) {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		commands.DownloadCmd.Usage()
		os.Exit(0)
	}
	//cmd.RunDownload(args)
}

func uploadCommand(args []string) {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		commands.UploadCmd.Usage()
		os.Exit(0)
	}
	//cmd.RunUpload(args)
}
