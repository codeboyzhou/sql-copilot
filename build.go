package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const TestCoverageDir = "coverage"

func main() {
	if len(os.Args) < 2 {
		showHelp()
	}

	task := os.Args[1]
	switch task {
	case "test":
		runTest(len(os.Args) == 3 && os.Args[2] == "-html")
	default:
		showHelp()
	}
}

func showHelp() {
	fmt.Println()
	fmt.Println("Usage: go run build.go <task>")
	fmt.Println()
	fmt.Println("Tasks:")
	fmt.Println("    test - Run unit tests")
	fmt.Println()
	os.Exit(0)
}

func run(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	fmt.Print(string(output))
	return nil
}

func runTest(openHtmlInBrowser bool) {
	if _, err := os.Stat(TestCoverageDir); os.IsExist(err) {
		os.RemoveAll(TestCoverageDir)
	}
	os.MkdirAll(TestCoverageDir, 0o755)

	coverprofile := fmt.Sprintf("%s/coverage.out", TestCoverageDir)
	if err := run("go", "test", "-coverprofile", coverprofile, "./slowsql"); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	coverhtml := fmt.Sprintf("%s/coverage.html", TestCoverageDir)
	run("go", "tool", "cover", "-html", coverprofile, "-o", coverhtml)
	fmt.Printf("✅ Test coverage report generated: %s\n", coverhtml)

	if openHtmlInBrowser {
		if runtime.GOOS == "windows" {
			run("cmd", "/c", "start", "", coverhtml)
		}
	} else {
		fmt.Println("💡 To view the coverage report, use: go run build.go test -html")
	}
}
