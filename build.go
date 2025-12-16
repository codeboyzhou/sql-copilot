//go:build build
// +build build

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const TestCoverageDir = "coverage"

const (
	EmojiSuccess = "✅"
	EmojiError   = "❌"
	EmojiTips    = "💡"
	EmojiRunning = "⏳"
)

const (
	GoFumptVersion   = "v0.9.2"
	GoImportsVersion = "v0.40.0"
	GolangciVersion  = "v2.7.2"
)

const (
	MySQLDockerImage        = "mysql:8.0.44-debian"
	DockerComposeTestFile   = "docker-compose.test.yml"
	MySQLTestContainerName  = "mysql-test"
	MySQLInitTimeoutSeconds = 10
	MySQLTestDSN            = "root:mysqlroot@tcp(127.0.0.1:3307)/test?parseTime=true"
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
	case "tools":
		runTools()
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
	case "ci":
		runIntegrationTest()
	default:
		showHelp()
	}
}

func showHelp() {
	fmt.Println()
	fmt.Println("Usage: go run build.go <task>")
	fmt.Println()
	fmt.Println("Tasks:")
	fmt.Println("    help      - Show this help information")
	fmt.Println("    tools     - Install required go tools")
	fmt.Println("    lint      - Run lint checks")
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

func runTools() {
	fmt.Printf("%s Installing gofumpt (%s)...\n", EmojiRunning, GoFumptVersion)
	gofumpt := fmt.Sprintf("mvdan.cc/gofumpt@%s", GoFumptVersion)
	if err := run("go", "install", gofumpt); err != nil {
		fmt.Printf("%s Error: failed to install gofumpt (%s)\n", EmojiError, GoFumptVersion)
		os.Exit(1)
	}
	fmt.Printf("%s Successfully installed gofumpt (%s)\n", EmojiSuccess, GoFumptVersion)

	fmt.Printf("%s Installing goimports (%s)...\n", EmojiRunning, GoImportsVersion)
	goimports := fmt.Sprintf("golang.org/x/tools/cmd/goimports@%s", GoImportsVersion)
	if err := run("go", "install", goimports); err != nil {
		fmt.Printf("%s Error: failed to install goimports (%s)\n", EmojiError, GoImportsVersion)
		os.Exit(1)
	}
	fmt.Printf("%s Successfully installed goimports (%s)\n", EmojiSuccess, GoImportsVersion)

	fmt.Printf("%s Installing golangci-lint (%s)...\n", EmojiRunning, GolangciVersion)
	golangci := fmt.Sprintf("github.com/golangci/golangci-lint/v2/cmd/golangci-lint@%s", GolangciVersion)
	if err := run("go", "install", golangci); err != nil {
		fmt.Printf("%s Error: failed to install golangci-lint (%s)\n", EmojiError, GolangciVersion)
		os.Exit(1)
	}
	fmt.Printf("%s Successfully installed golangci-lint (%s)\n", EmojiSuccess, GolangciVersion)
}

func runLint() {
	fmt.Printf("%s Running:  go mod tidy\n", EmojiRunning)
	// should never fail but just in case
	if err := run("go", "mod", "tidy"); err != nil {
		fmt.Printf("%s Error: go mod tidy failed\n", EmojiError)
		os.Exit(1)
	}
	fmt.Printf("%s Finished: go mod tidy\n", EmojiSuccess)

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
		if err := os.RemoveAll(TestCoverageDir); err != nil {
			fmt.Printf("%s Error: %v\n", EmojiError, err)
			os.Exit(1)
		}
	}

	if err := os.MkdirAll(TestCoverageDir, 0o755); err != nil {
		fmt.Printf("%s Error: %v\n", EmojiError, err)
		os.Exit(1)
	}

	coverprofile := fmt.Sprintf("%s/coverage.out", TestCoverageDir)
	if err := run("go", "test", "-v", "-coverprofile", coverprofile, "./..."); err != nil {
		fmt.Printf("%s Error: there are test failures\n", EmojiError)
		os.Exit(1)
	}

	coverhtml := fmt.Sprintf("%s/coverage.html", TestCoverageDir)
	if err := run("go", "tool", "cover", "-html", coverprofile, "-o", coverhtml); err != nil {
		fmt.Printf("%s Error: %v\n", EmojiError, err)
		os.Exit(1)
	}
	fmt.Printf("%s Test coverage report generated: %s\n", EmojiSuccess, coverhtml)

	if openHtmlInBrowser {
		if runtime.GOOS == "windows" {
			if err := run("cmd", "/c", "start", "", coverhtml); err != nil {
				fmt.Printf("%s Error: %v\n", EmojiError, err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Printf("%s To view the coverage report, use: go run build.go test -html\n", EmojiTips)
	}
}

func runBenchmark(bench string) {
	fmt.Printf("%s Running benchmark...It may take a while\n", EmojiRunning)
	if err := run("go", "test", "-bench", bench, "-benchtime", "5s", "./..."); err != nil {
		fmt.Printf("%s Error: %v\n", EmojiError, err)
		os.Exit(1)
	}
	fmt.Printf("%s Benchmark completed successfully\n", EmojiSuccess)
}

func runIntegrationTest() {
	fmt.Printf("%s Running integration tests...\n", EmojiRunning)

	fmt.Printf("%s Pulling Docker image %s for MySQL test container...\n", EmojiRunning, MySQLDockerImage)
	if err := run("docker", "pull", MySQLDockerImage); err != nil {
		fmt.Printf("%s Error: %v\n", EmojiError, err)
		os.Exit(1)
	}
	fmt.Printf("%s Successfully pulled Docker image %s for MySQL test container\n", EmojiSuccess, MySQLDockerImage)

	fmt.Printf("%s Running MySQL test container...\n", EmojiRunning)
	if err := run("docker-compose", "-f", DockerComposeTestFile, "up", "-d"); err != nil {
		fmt.Printf("%s Error: %v\n", EmojiError, err)
		os.Exit(1)
	}
	fmt.Printf("%s Successfully started MySQL test container: %s\n", EmojiSuccess, MySQLTestContainerName)

	fmt.Printf("%s Waiting %d seconds for MySQL to initialize...\n", EmojiRunning, MySQLInitTimeoutSeconds)
	time.Sleep(MySQLInitTimeoutSeconds * time.Second)

	fmt.Printf("%s Setting TEST_DB_DSN environment variable...\n", EmojiRunning)
	os.Setenv("TEST_DB_DSN", MySQLTestDSN)
	fmt.Printf("%s TEST_DB_DSN set to: %s\n", EmojiSuccess, os.Getenv("TEST_DB_DSN"))

	fmt.Printf("%s Running Go integration tests...\n", EmojiRunning)
	if err := run("go", "test", "-v", "-tags=integration", "./integration_test/"); err != nil {
		fmt.Printf("%s Error: %v\n", EmojiError, err)
		os.Exit(1)
	}
	fmt.Printf("%s Successfully completed integration tests\n", EmojiSuccess)

	fmt.Printf("%s Stopping and removing MySQL test container...\n", EmojiRunning)
	if err := run("docker-compose", "-f", DockerComposeTestFile, "down"); err != nil {
		fmt.Printf("%s Error: %v\n", EmojiError, err)
		os.Exit(1)
	}
	fmt.Printf("%s Successfully stopped and removed MySQL test container: %s\n", EmojiSuccess, MySQLTestContainerName)
}
