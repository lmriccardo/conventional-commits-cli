package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitInfo struct {
	Reponame    string   // The name of the current repository
	Branches    []string // All the branches for the current repository
	Remotes     []string // All remotes for the current repository
	Curr_branch string   // The current branch name
	Curr_remote string   // The remote for the current branch
	Commit_str  string   // The commit message string
}

func extractRepoName(url string) string {
	// Remove protocol (http://, https://, git@)
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "git@")
	url = strings.TrimSuffix(url, ".git")

	// Remove any possible :
	parts := strings.Split(url, ":")
	if len(parts) > 1 {
		url = parts[len(parts)-1]
	}

	// Split the string by the slash and take the last part
	parts = strings.Split(url, "/")
	sig_parts := []string{parts[len(parts)-2], parts[len(parts)-1]}
	return strings.Join(sig_parts, "/")
}

// Returns the name of the repository
func getRepositoryName(gitdir string) (string, error) {
	// The name of the repository can be found in the config file
	// Format the file path and reads its content, if any error occurs
	// then it will returns an empty string as long with the error
	config_file := filepath.Join(gitdir, "config")
	data, err := os.ReadFile(config_file)
	if err != nil {
		return "", err
	}

	// Parse the config file to find the repository name.
	// Usually, it is indicated into the remote url
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "url =") {
			continue
		}

		url := strings.TrimSpace(strings.Split(line, "=")[1])
		return extractRepoName(url), nil
	}

	return "", err
}

// Returns the name of all branches
func getAllBranches(rootdir string, level int) ([]string, error) {
	// This function recursively dive into folders to get the
	// name of the file contained. Each file represents a branch
	branches := make([]string, 0)
	entries, err := os.ReadDir(rootdir)
	if err != nil {
		return branches, err
	}

	// Loop for all the entries of the folder and recurse if necessary
	for _, entry := range entries {
		if !entry.IsDir() {
			branch_name := entry.Name()

			if level >= 1 {
				parts := strings.Split(rootdir, "/")
				end_idx := len(parts)
				suffix := strings.Join(parts[end_idx-level:], "/")
				branch_name = suffix + "/" + branch_name
			}

			branches = append(branches, branch_name)
			continue
		}

		new_root := filepath.Join(rootdir, entry.Name())
		sub_branches, err := getAllBranches(new_root, level+1)
		if err != nil {
			return nil, err
		}

		// Otherwise append the sub branches to the list
		branches = append(branches, sub_branches...)
	}

	return branches, nil
}

// Returns the current branch
func getCurrentBranch(gitdir string) (string, error) {
	// The current branch name should be in HEAD file
	head_file := filepath.Join(gitdir, "HEAD")
	data, err := os.ReadFile(head_file)
	if err != nil {
		return "", nil
	}

	// Usually there should be only one row (not too sure about that)
	line := strings.Split(string(data), "\n")[0]
	parts := strings.Split(line, "/")

	// Search the index of the "heads" string in the array of parts
	head_idx := -1
	for idx := range parts {
		if strings.Contains(parts[idx], "heads") {
			head_idx = idx
			break
		}
	}

	return strings.Join(parts[head_idx+1:], "/"), nil
}

func getAllRemotes(gitdir string) ([]string, error) {
	config_file := filepath.Join(gitdir, "config")
	data, err := os.ReadFile(config_file)
	if err != nil {
		return nil, err
	}

	// Parse the config file to find the repository name.
	// Usually, it is indicated into the remote url
	lines := strings.Split(string(data), "\n")
	remotes := make([]string, 0)
	for _, line := range lines {
		if !strings.HasPrefix(line, "[remote ") {
			continue
		}

		start_idx := len("[remote ")
		end_idx := len(line) - 1
		remote := line[start_idx+1 : end_idx-1]
		remotes = append(remotes, remote)
	}

	return remotes, nil
}

func GetGitInfo(rootpath string) *GitInfo {
	// First of all, we need to chech that the current folder
	// is a git repository, meaning the .git folder exists
	git_dir := filepath.Join(rootpath, ".git")
	if _, err := os.Stat(git_dir); os.IsNotExist(err) {
		return nil
	}

	// Initialize the return value
	gitinfo := new(GitInfo)

	// Get the repository name
	repo_name, err := getRepositoryName(git_dir)
	if err != nil {
		return nil
	}

	gitinfo.Reponame = repo_name // Set the repository name

	// Get all the branches from the .git/refs/head/ folder
	branches, err := getAllBranches(filepath.Join(git_dir, "refs", "heads"), 0)
	if err != nil {
		return nil
	}

	gitinfo.Branches = branches // Set all the branches name

	// Get the current branch name
	branch_name, err := getCurrentBranch(git_dir)
	if err != nil {
		return nil
	}

	gitinfo.Curr_branch = branch_name // Set the branch name

	// Get all remotes
	remotes, err := getAllRemotes(git_dir)
	if err != nil {
		return nil
	}

	gitinfo.Remotes = remotes // Set the remotes to the info structure

	return gitinfo
}

func (gi *GitInfo) FinalizeCommit() {
	// Print some useful informations
	gitstatus := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	gitstatus.Stdout = &out
	gitstatus.Run()

	status := strings.TrimSpace(out.String())
	if len(status) < 1 {
		fmt.Println("[*] No changes to commit. Exiting ...")
		os.Exit(1)
	}

	fmt.Println("[*] Showing the current status")
	fmt.Println()
	fmt.Printf("%s\n", status)
	fmt.Println()
	fmt.Println("[*] Previous changes needs to be staged before commiting.")
	fmt.Println("[*] Running commands: <git add .> and <git commit -m ...> (Press ENTER to run, CTRL + C for exit)")
	fmt.Scanln()

	// Run Git add command
	gitadd := exec.Command("git", "add", ".")
	gitadd.Stderr = os.Stderr
	gitadd.Stdout = os.Stdout
	err := gitadd.Run()
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		return
	}

	// Run git commit
	gitcommit := exec.Command("git", "commit", "-m", gi.Commit_str)
	gitcommit.Stderr = os.Stderr
	gitadd.Stdout = os.Stdout
	err = gitcommit.Run()
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		return
	}

	fmt.Println("\n[*] Pushing changes into remote. (Press ENTER to run, CTRL + C for exit)")
	fmt.Scanln()

	// Run git push
	gitpush := exec.Command("git", "push", "--set-upstream", gi.Curr_remote, gi.Curr_branch)
	gitpush.Stderr = os.Stderr
	gitpush.Stdout = os.Stdout
	err = gitpush.Run()
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		return
	}
}
