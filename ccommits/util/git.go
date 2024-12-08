package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const separator string = string(os.PathSeparator)

type GitInfo struct {
	Reponame    string   // The name of the current repository
	Branches    []string // All the branches for the current repository
	Remotes     []string // All remotes for the current repository
	Curr_branch string   // The current branch name
	Curr_remote string   // The remote for the current branch
	Commit_str  string   // The commit message string
	PrevContent string   // The previous content of the .git file (only for worktrees)
	GitDir      string   // The root folder of git
	TargetPath  string   // The target path of all git commands
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
				parts := strings.Split(rootdir, separator)
				end_idx := len(parts)
				suffix := strings.Join(parts[end_idx-level:], separator)
				branch_name = suffix + separator + branch_name
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

func getGitInfo(rootpath, srcpath, entrypath string) *GitInfo {
	// Initialize the return value
	gitinfo := new(GitInfo)
	gitinfo.PrevContent = ""
	gitinfo.TargetPath = rootpath

	// First of all, we need to chech that the current folder
	// is a git repository, meaning the .git folder exists
	git_dir := filepath.Join(rootpath, ".git")
	info, err := os.Stat(git_dir)
	if os.IsNotExist(err) {
		fmt.Printf("An Error occurred: %s\n", err)
		return nil
	}

	gitinfo.GitDir = git_dir // Set the git folder inside the structure
	branch_dir := git_dir    // We need to set also the directory where to take the branch

	// In case of worktrees, we need to check whether the git_dir is actually
	// a folder or a file linking to the real repository folder
	if !info.IsDir() {
		// Then we can read the content of the file and retrieve the actual folder
		data, _ := os.ReadFile(git_dir)
		data_str := string(data)
		data_str = data_str[0 : len(data_str)-1]
		parts := strings.Split(data_str, ": ")
		branch_dir = parts[len(parts)-1] // Take the branch folder

		// If the entry path of the container, meaning the one the user have
		// previously bind mounted and set as working folder is different
		// from the source folder, i.e., the path of the mount in the host
		// filesystem, we need to change the absolute path to the branch folder.
		// Moreover, it is necessary to change the content of the .git file otherwise
		// git will not able to work properly.
		if strings.Compare(srcpath, entrypath) != 0 {
			parts = strings.Split(branch_dir, srcpath)
			branch_dir = filepath.Join(entrypath, parts[1])

			// Open the .git file in write mode
			file, _ := os.OpenFile(git_dir, os.O_TRUNC|os.O_WRONLY, 0644)
			gitinfo.PrevContent = data_str
			newcontent := fmt.Sprintf("gitdir: %s\n", branch_dir)
			file.WriteString(newcontent)
		}

		config_cmd := exec.Command("git", "config", "--add", "safe.directory", rootpath)
		err = config_cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}

		// The master .git folder is located relatively to the current
		// worktree branch folder at <curr_path>/<content-commondir>
		// where commondir is a text file containing the relative path
		// to the original git repository folder
		tmp_path := filepath.Join(branch_dir, "commondir")
		common_dir, _ := os.ReadFile(tmp_path)
		common_dir_str := string(common_dir)
		common_dir_str = common_dir_str[0 : len(common_dir_str)-1]
		git_dir = strings.Join([]string{branch_dir, common_dir_str}, separator)
	}

	// Get the repository name
	repo_name, err := getRepositoryName(git_dir)
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		return nil
	}

	gitinfo.Reponame = repo_name // Set the repository name

	// Get all the branches from the .git/refs/head/ folder
	branches, err := getAllBranches(filepath.Join(git_dir, "refs", "heads"), 0)
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		return nil
	}

	gitinfo.Branches = branches // Set all the branches name

	// Get the current branch name
	branch_name, err := getCurrentBranch(branch_dir)
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		return nil
	}

	gitinfo.Curr_branch = branch_name // Set the branch name

	// Get all remotes
	remotes, err := getAllRemotes(git_dir)
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		return nil
	}

	gitinfo.Remotes = remotes // Set the remotes to the info structure

	return gitinfo
}

func checkChangesToCommit(gi *GitInfo) {
	// Get the status of the current branch. We need to check if there
	// are changes that needs to be committed
	gitstatus := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	gitstatus.Stdout = &out
	gitstatus.Run()

	status := strings.TrimSpace(out.String())
	if len(status) < 1 {
		fmt.Println("[*] No changes to commit. Exiting ...")
		gi.RestorePreviousContent()
		os.Exit(1)
	}

	fmt.Println("[*] Showing the current status")
	fmt.Println()
	fmt.Printf("%s\n", status)
	fmt.Println()
}

func GetGitRepositoryInformation(remote_name, tgfolder, srcfolder, entrypath string) *GitInfo {
	fmt.Println("------------------------- GIT REPOSITORY GATHERING -------------------------")

	gitinfo := getGitInfo(tgfolder, srcfolder, entrypath)
	if gitinfo == nil {
		fmt.Println("Current folder does not belong to a repository. Exiting ...")
		os.Exit(1)
	}

	fmt.Printf("DETECTED REPOSITORY: \033[3m%s\033[0m\n", gitinfo.Reponame)
	fmt.Printf("DETECTED CURRENT BRANCH: \033[3m%s\033[0m\n", gitinfo.Curr_branch)
	fmt.Printf("DETECTED REPOSITORY BRANCHES: \033[3m%s\033[0m\n", strings.Join(gitinfo.Branches, ", "))
	fmt.Printf("DETECTED POSSIBLE REMOTES: \033[3m%s\033[0m\n", strings.Join(gitinfo.Remotes, ", "))

	if len(remote_name) < 1 && len(gitinfo.Remotes) > 1 {
		fmt.Print("\n[*] Please Choose a remote: ")
		fmt.Scanln(&gitinfo.Curr_remote)
	} else if len(remote_name) > 1 {
		gitinfo.Curr_remote = remote_name
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

	// We need to change the cwd to the target path
	cwd, _ := os.Getwd()
	if strings.Compare(cwd, gitinfo.TargetPath) != 0 {
		os.Chdir(gitinfo.TargetPath)
	}

	// Check for changes to be committed
	checkChangesToCommit(gitinfo)

	fmt.Println("----------------------------------------------------------------------------")

	return gitinfo
}

func (gi *GitInfo) FinalizeCommit(flag bool) {
	fmt.Println("[*] Previous changes needs to be staged before commiting.")

	if !flag {
		fmt.Println("[*] Running commands: <git add .> and <git commit -m ...> (Press ENTER to run, CTRL + C for exit)")
		fmt.Scanln()
	} else {
		fmt.Println("[*] Running commands: <git add .> and <git commit -m ...>")
	}

	// Run Git add command
	gitadd := exec.Command("git", "add", ".")
	gitadd.Stderr = os.Stderr
	gitadd.Stdout = os.Stdout
	err := gitadd.Run()
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		gi.RestorePreviousContent()
		return
	}

	// Run git commit
	gitcommit := exec.Command("git", "commit", "-m", gi.Commit_str)
	gitcommit.Stderr = os.Stderr
	gitadd.Stdout = os.Stdout
	err = gitcommit.Run()
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		gi.RestorePreviousContent()
		return
	}

	if !flag {
		fmt.Println("\n[*] Pushing changes into remote. (Press ENTER to run, CTRL + C for exit)")
		fmt.Scanln()
	} else {
		fmt.Println("\n[*] Pushing changes into remote.")
	}

	// Run git push
	gitpush := exec.Command("git", "push", "--set-upstream", gi.Curr_remote, gi.Curr_branch)
	gitpush.Stderr = os.Stderr
	gitpush.Stdout = os.Stdout
	err = gitpush.Run()
	if err != nil {
		fmt.Printf("An Error occurred: %s\n", err)
		gi.RestorePreviousContent()
		return
	}

	gi.RestorePreviousContent()
}

func (gi *GitInfo) RestorePreviousContent() {
	// Restore the previous state of the .git file (if necessary)
	if len(gi.PrevContent) > 0 {
		file, _ := os.OpenFile(gi.GitDir, os.O_WRONLY, 0644)
		file.WriteString(gi.PrevContent + "\n")
	}
}
