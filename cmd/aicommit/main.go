package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/loissascha/aicommit/internal/ai"
)

func envFileExists() bool {
	_, err := os.Stat(".env")
	if err == nil {
		return true
	}
	return false
}

func main() {
	showTokens := flag.Bool("tokens", false, "add this flag to print out the used tokens.")
	confirmFlag := flag.Bool("confirm", false, "require confirmation.")
	flag.Parse()

	if envFileExists() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error reading .env file")
		}
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set in .env file or environment variable")
	}

	cmd := exec.Command("git", "diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalln("Error getting git diff of staged files. Is this a git repo?")
	}

	if strings.TrimSpace(string(out)) == "" {
		log.Fatalln("There are no staged files.")
	}

	header, message, err := ai.GenerateCommitMessage(string(out), *showTokens)
	retriesLeft := 2
	if err != nil {
		for retriesLeft > 0 {
			retriesLeft--
			header, message, err = ai.GenerateCommitMessage(string(out), *showTokens)
			if err == nil {
				retriesLeft = 0
			}
		}
	}
	if err != nil {
		log.Fatalln("Error getting AI commit message:", err.Error())
	}

	if *confirmFlag {
		fmt.Println("")
		fmt.Println(header)
		fmt.Println("---------------")
		fmt.Println(message)
		fmt.Println("")

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
		log.Fatalln("Error commiting:", err.Error())
	}
}
