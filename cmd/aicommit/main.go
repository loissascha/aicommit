package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/loissascha/aicommit/internal/ai"
)

func main() {
	showTokens := flag.Bool("tokens", false, "add this flag to print out the used tokens.")
	confirmFlag := flag.Bool("confirm", false, "require confirmation.")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("GEMINI_API_KEY not set in .env file or environment variable")
	}

	cmd := exec.Command("git", "diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(string(out)) == "" {
		fmt.Println("There are no staged files.")
		return
	}

	header, message, err := ai.GenerateCommitMessage(string(out), *showTokens)
	if err != nil {
		panic(err)
	}

	if *confirmFlag {
		fmt.Println(header)
		fmt.Println("---------------")
		fmt.Println(message)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Would you like to accept this commit? [yes/no] (default:yes): ")
		conf, _ := reader.ReadString('\n')
		conf = strings.TrimSpace(conf)
		if conf == "" || conf == "yes" {
			fmt.Println("")
		} else {
			return
		}
	}

	cmd = exec.Command("git", "commit", "-m", header, "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
