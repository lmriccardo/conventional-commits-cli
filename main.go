package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"example.com/ccommits/ccommits"
	"example.com/ccommits/ccommits/util"
)

func main() {
	fmt.Println(ccommits.TITLE)
	fmt.Printf("Running Version: %s\n", ccommits.VERSION)
	fmt.Println("GitHub Repository: https://github.com/lmriccardo/conventional-commits-cli.git")
	fmt.Println()

	cwd, _ := os.Getwd()
	gitinfo := util.GetGitInfo(cwd)
	if gitinfo == nil {
		fmt.Println("Current folder does not belong to a repository. Exiting ...")
		os.Exit(1)
	}

	// Define the expected input command line argument
	remote_name := flag.String("remote", "", "The chosen remote name")
	flag.Parse()

	fmt.Printf("DETECTED REPOSITORY: \033[3m%s\033[0m\n", gitinfo.Reponame)
	fmt.Printf("DETECTED CURRENT BRANCH: \033[3m%s\033[0m\n", gitinfo.Curr_branch)
	fmt.Printf("DETECTED REPOSITORY BRANCHES: \033[3m%s\033[0m\n",
		strings.Join(gitinfo.Branches, ", "))

	fmt.Printf("DETECTED POSSIBLE REMOTES: \033[3m%s\033[0m\n",
		strings.Join(gitinfo.Remotes, ", "))

	if len(*remote_name) < 1 && len(gitinfo.Branches) > 1 {
		fmt.Print("\n[*] Please Choose a remote: ")
		fmt.Scanln(&gitinfo.Curr_remote)
	} else if len(*remote_name) > 1 {
		gitinfo.Curr_remote = *remote_name
	} else if len(gitinfo.Remotes) == 1 {
		gitinfo.Curr_remote = gitinfo.Remotes[0]
	}

	// Check that at least a name has been given
	if len(gitinfo.Curr_remote) < 1 {
		fmt.Println("A remote name must be choosen. Exiting ...")
		os.Exit(1)
	}

	// Check that the remote name is inside the list of all remotes
	result := false
	for _, remote := range gitinfo.Remotes {
		if strings.Compare(gitinfo.Curr_remote, remote) == 0 {
			result = true
		}
	}

	if !result {
		fmt.Println("A valid remote name must be given. Exiting ...")
		os.Exit(1)
	}

	fmt.Println("\n[*] Running conventional commits cli app")
	time.Sleep(time.Second)

	app := ccommits.CCommitWindow_new(gitinfo)
	fmt_commit := app.Run()
	if len(fmt_commit) < 1 {
		fmt.Println("Invalid formatted conventional commit. Exiting ...")
		os.Exit(1)
	}

	fmt.Println("[*] Following result obtained")
	fmt.Println()
	fmt.Printf("%s\n\n", fmt_commit)
	fmt.Println()

	fmt.Println("[*] Finalizing the Commit and Closing")
	gitinfo.Commit_str = fmt_commit
	gitinfo.FinalizeCommit()
}
