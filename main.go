package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lmriccardo/conventional-commits-cli/ccommits"
	"github.com/lmriccardo/conventional-commits-cli/ccommits/util"
)

func main() {
	fmt.Println(ccommits.TITLE)
	fmt.Printf("Running Version: %s\n", ccommits.VERSION)
	fmt.Println("GitHub Repository: https://github.com/lmriccardo/conventional-commits-cli.git")
	fmt.Println()

	cwd, _ := os.Getwd()

	// Before getting git repository info it must check if the current environment
	// is a docker container by using the defined heuristics
	fmt.Println("------------------------- CONTAINER CHECKS ---------------------------------")
	fmt.Print("[*] Is the current environment under a container? ")
	is_container, container_id, _ := util.IsContainerEnvironment()
	if !is_container {
		fmt.Println("No")
	} else {
		fmt.Printf("Yes.\nDETECTED CONTAINER ID: %s\n", container_id)

		// Check if the current working folder is mounted inside the container
		fmt.Print("\n[*] Checking if the CWD is bind mounted: ")
		fs_output := util.IsContainerFolderMounted(cwd)

		if fs_output == nil {
			fmt.Println("No")
		} else {
			fs_data := fs_output.Filesystems[0] // Select the only entry
			fmt.Println("Yes")
			fmt.Printf("   HOST SOURCE: %s\n", fs_data.Source)
			fmt.Printf("   CONTAINER TARGET: %s\n", fs_data.Target)
		}

		// Get the ccommits target folder from the environment variable
		fmt.Print("\n[*] Is the CCOMMITS_WD env variable set: ")
		target_folder, _ := util.GetContainerEnvironmentVariable("CCOMMITS_WD")
		if len(target_folder) < 1 {
			// If the environment variable is not set, then ask the user
			fmt.Println("No")
			fmt.Print("[*] Enter the (relative) target folder (Leave blank for CWD): ")
			fmt.Scanln(&target_folder)
			if len(target_folder) < 1 {
				target_folder = cwd // Set to the current working folder
			}
		}

		fmt.Printf("SELECTED TARGET FOLDER: %s\n", target_folder)
	}

	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("------------------------- GIT REPOSITORY GATHERING -------------------------")

	gitinfo := util.GetGitInfo(cwd)
	if gitinfo == nil {
		fmt.Println("Current folder does not belong to a repository. Exiting ...")
		os.Exit(1)
	}

	// Define the expected input command line argument
	remote_name := flag.String("remote", "", "The chosen remote name")
	yes_flag := flag.Bool("yes", false, "Skip all user input pauses when finalizing commit")
	flag.Parse()

	fmt.Printf("DETECTED REPOSITORY: \033[3m%s\033[0m\n", gitinfo.Reponame)
	fmt.Printf("DETECTED CURRENT BRANCH: \033[3m%s\033[0m\n", gitinfo.Curr_branch)
	fmt.Printf("DETECTED REPOSITORY BRANCHES: \033[3m%s\033[0m\n",
		strings.Join(gitinfo.Branches, ", "))

	fmt.Printf("DETECTED POSSIBLE REMOTES: \033[3m%s\033[0m\n",
		strings.Join(gitinfo.Remotes, ", "))

	if len(*remote_name) < 1 && len(gitinfo.Remotes) > 1 {
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

	fmt.Println("----------------------------------------------------------------------------")

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

	fmt.Println("------------------------- FINALIZING THE COMMIT ----------------------------")

	fmt.Println("[*] Finalizing the Commit and Closing")
	gitinfo.Commit_str = fmt_commit
	gitinfo.FinalizeCommit(*yes_flag)

	fmt.Println("----------------------------------------------------------------------------")
}
