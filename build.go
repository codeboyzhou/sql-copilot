package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const TestCoverageDir = "coverage"

const (
	EmojiSuccess = "✅"
	EmojiError   = "❌"
	EmojiTips    = "💡"
	EmojiRunning = "⏳"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	task := os.Args[1]
	switch task {
	case "help":
		showHelp()
		os.Exit(0)
	case "lint":
		runLint()
	case "test":
		runTest(len(os.Args) == 3 && os.Args[2] == "-html")
	case "benchmark":
		if len(os.Args) < 3 {
			fmt.Printf("%s Error: Please specify the benchmark to run\n", EmojiError)
			fmt.Printf("%s Example: go run build.go benchmark BenchmarkParseSlowLog\n", EmojiTips)
			fmt.Printf("%s To run all benchmarks, use: go run build.go benchmark .\n", EmojiTips)
			os.Exit(1)
		}
		bench := os.Args[2]
		runBenchmark(bench)
	default:
		showHelp()
	}
}

func showHelp() {
	fmt.Println()
	fmt.Println("Usage: go run build.go <task>")
	fmt.Println()
	fmt.Println("Tasks:")
	fmt.Println("    test      - Run unit tests")
	fmt.Println("    benchmark - Run benchmarks")
	fmt.Println()
}

func run(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()

	if len(output) > 0 {
		fmt.Print(string(output))
	}

	if err != nil {
		return err
	}

	return nil
}

func runLint() {
	fmt.Printf("%s Running:  go fmt\n", EmojiRunning)
	// should never fail but just in case
	if err := run("go", "fmt"); err != nil {
		fmt.Printf("%s Error: go fmt failed\n", EmojiError)
		os.Exit(1)
	}
	fmt.Printf("%s Finished: go fmt\n", EmojiSuccess)

	fmt.Printf("%s Running:  gofumpt -w .\n", EmojiRunning)
	// should never fail but just in case
	if err := run("gofumpt", "-w", "."); err != nil {
		fmt.Printf("%s Error: gofumpt -w . failed\n", EmojiError)
		os.Exit(1)
	}
	fmt.Printf("%s Finished: gofumpt -w .\n", EmojiSuccess)

	fmt.Printf("%s Running:  golangci-lint run ./... --config .golangci.yml\n", EmojiRunning)
	if err := run("golangci-lint", "run", "./...", "--config", ".golangci.yml"); err != nil {
		fmt.Printf("%s Error: golangci-lint failed\n", EmojiError)
		os.Exit(1)
	}
	fmt.Printf("%s Finished: golangci-lint run ./... --config .golangci.yml\n", EmojiSuccess)
}

func runTest(openHtmlInBrowser bool) {
	if _, err := os.Stat(TestCoverageDir); os.IsExist(err) {
		os.RemoveAll(TestCoverageDir)
	}
	os.MkdirAll(TestCoverageDir, 0o755)

	coverprofile := fmt.Sprintf("%s/coverage.out", TestCoverageDir)
	if err := run("go", "test", "-coverprofile", coverprofile, "./slowsql"); err != nil {
		fmt.Printf("%s Error: %s\n", EmojiError, err)
		os.Exit(1)
	}

	coverhtml := fmt.Sprintf("%s/coverage.html", TestCoverageDir)
	run("go", "tool", "cover", "-html", coverprofile, "-o", coverhtml)
	fmt.Printf("%s Test coverage report generated: %s\n", EmojiSuccess, coverhtml)

	if openHtmlInBrowser {
		if runtime.GOOS == "windows" {
			run("cmd", "/c", "start", "", coverhtml)
		} else {
			fmt.Printf("%s To view the coverage report, use: go run build.go test -html\n", EmojiTips)
		}
	}
}

func runBenchmark(bench string) {
	fmt.Printf("%s Running benchmark...It may take a while\n", EmojiRunning)
	if err := run("go", "test", "-bench", bench, "-benchtime", "5s", "./slowsql"); err != nil {
		fmt.Printf("%s Error: %s\n", EmojiError, err)
		os.Exit(1)
	}
	fmt.Printf("%s Benchmark completed successfully\n", EmojiSuccess)
}
