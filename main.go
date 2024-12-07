package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lmriccardo/conventional-commits-cli/ccommits"
	"github.com/lmriccardo/conventional-commits-cli/ccommits/util"
)

func main() {
	fmt.Println(ccommits.TITLE)
	fmt.Printf("Running Version: %s\n", ccommits.VERSION)
	fmt.Println("GitHub Repository: https://github.com/lmriccardo/conventional-commits-cli.git")
	fmt.Println()

	// Define the expected input command line argument
	remote_name := flag.String("remote", "", "The chosen remote name")
	yes_flag := flag.Bool("yes", false, "Skip all user input pauses when finalizing commit")
	flag.Parse()

	cwd, _ := os.Getwd()

	// Before getting git repository info it must check if the current environment
	// is a docker container by using the defined heuristics
	target_folder, src_folder, entry_path := util.PerformContainerChecks(cwd)

	// Gets repository information
	gitinfo := util.GetGitRepositoryInformation(*remote_name, target_folder, src_folder, entry_path)

	fmt.Println("\n[*] Running conventional commits cli app")
	time.Sleep(time.Second)

	app := ccommits.CCommitWindow_new(gitinfo)
	fmt_commit := app.Run()
	if len(fmt_commit) < 1 {
		fmt.Println("Invalid formatted conventional commit. Exiting ...")
		gitinfo.RestorePreviousContent()
		os.Exit(1)
	}

	fmt.Println("[*] Following result obtained")
	fmt.Println()
	fmt.Printf("%s\n\n", fmt_commit)
	fmt.Println()

	fmt.Println("------------------------- FINALIZING THE COMMIT ----------------------------")

	fmt.Println("[*] Finalizing the Commit and Closing")
	gitinfo.Commit_str = fmt_commit
	gitinfo.FinalizeCommit(*yes_flag)

	fmt.Println("----------------------------------------------------------------------------")
}
