package main

import (
	"fmt"
	"os"
	"strings"

	"example.com/ccommits/src"
	"example.com/ccommits/src/util"
)

func main() {
	fmt.Println(src.TITLE)
	fmt.Printf("Running Version: %s\n", src.VERSION)
	fmt.Println("GitHub Repository: https://github.com/lmriccardo/conventional-commits-cli.git")
	fmt.Println()

	cwd, _ := os.Getwd()
	gitinfo := util.GetGitInfo(cwd)
	if gitinfo == nil {
		fmt.Println("Current folder does not belong to a repository. Exiting ...")
		os.Exit(1)
	}

	fmt.Printf("[*] Detected Repository Named: \033[3m%s\033[0m\n", gitinfo.Reponame)
	fmt.Printf("[*] Detected Current Branch: \033[3m%s\033[0m\n", gitinfo.Curr_branch)
	fmt.Printf("[*] Detected Repository Branches: \033[3m%s\033[0m\n",
		strings.Join(gitinfo.Branches, ", "))

	fmt.Println("[*] Running conventional commits cli app (Press ENTER to continue)")
	fmt.Scanln()

	app := src.CCommitWindow_new(gitinfo)
	fmt_commit := app.Run()
	if len(fmt_commit) < 1 {
		fmt.Println("Invalid formatted conventional commit. Exiting ...")
		os.Exit(1)
	}

	fmt.Println("[*] Following result obtained")
	fmt.Println()
	fmt.Printf("%s\n\n", fmt_commit)
}
