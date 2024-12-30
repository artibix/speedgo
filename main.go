package main

import (
	"context"
	"fmt"
	"os"
	"speedgo/commands"
	"speedgo/core"
)

func main() {
	// Create a base context that can be used across the application
	ctx := context.Background()

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "ping", "p":
		if err := pingCommand(ctx, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "download", "d":
		if err := downloadCommand(ctx, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "upload", "u":
		if err := uploadCommand(ctx, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
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
	fmt.Println("\nCommands:")
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

func pingCommand(ctx context.Context, args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		commands.PingCmd.Usage()
		return nil
	}
	return core.RunPing(ctx, args)
}

func downloadCommand(ctx context.Context, args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		commands.DownloadCmd.Usage()
		return nil
	}
	return core.RunDownload(ctx, args)
}

func uploadCommand(ctx context.Context, args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		commands.UploadCmd.Usage()
		return nil
	}
	return core.RunUpload(ctx, args)
}
