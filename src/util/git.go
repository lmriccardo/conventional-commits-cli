package util

import (
	"os"
	"path/filepath"
	"strings"
)

type GitInfo struct {
	Reponame    string   // The name of the current repository
	Branches    []string // All the branches for the current repository
	Curr_branch string   // The current branch name
	Curr_remote string   // The remote for the current branch
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
func getAllBranches(rootdir string) ([]string, error) {
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
			branches = append(branches, entry.Name())
			continue
		}

		new_root := filepath.Join(rootdir, entry.Name())
		sub_branches, err := getAllBranches(new_root)
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
	return parts[len(parts)-1], nil
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
	branches, err := getAllBranches(filepath.Join(git_dir, "refs", "heads"))
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

	return gitinfo
}
