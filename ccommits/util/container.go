package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Filesystem struct {
	Target  string `json:"target"`
	Source  string `json:"source"`
	Fstype  string `json:"fstype"`
	Options string `json:"options"`
}

type FindmntOutput struct {
	Filesystems []Filesystem `json:"filesystems"`
}

// Take all the environment variables and search for the one
// given as input to the function
func GetContainerEnvironmentVariable(varname string) (string, error) {
	variables := os.Environ()
	for _, element := range variables {
		parts := strings.Split(element, "=")
		if strings.Compare(parts[0], varname) == 0 {
			return parts[1], nil
		}
	}

	// If no variable has been found it returns an error
	return "", errors.New("no environment variable matches the provided key")
}

// Read the input cgroup file and returns the container id (in case)
func scanCgroupFile(filepath string) (string, error) {
	// Check that the file exists and can be opened
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err // Check for error and in case returns it
	}

	// First we need to check that the content of the file matches
	// what we expect to find in a typical docker container
	pattern := `^\d+:[a-zA-Z_-]*(=[a-zA-Z_-]+)?:/docker/[a-zA-Z0-9]+$`
	re := regexp.MustCompile(pattern) // Compile the pattern or error
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Check if the current line matches the pattern
		if re.MatchString(line) {
			docker_part := strings.Split(line, ":")[2]
			container_id := strings.Split(docker_part, "/")[2]
			return container_id, nil
		}
	}

	return "", errors.New("no container detected")
}

// Check if the current running environment is inside a docker container
func IsContainerEnvironment() (bool, string, error) {
	// There are some heuristic that can be used to identify
	// if the current running environment is a container or not.
	cgroups_file := "/proc/self/cgroup"
	dockerenv_file := "/.dockerenv"

	// Usually, if ccommits is started using the default command and
	// nothing is changed about cgroups, then the container will use
	// the default one created by docker.
	container_id, err := scanCgroupFile(cgroups_file)
	if err != nil {
		return false, "", err // Returns the error and false
	}

	// Finally, check if the typical .dockerenv file created
	// at the root of the container filesystem exists
	if _, err := os.Stat(dockerenv_file); os.IsNotExist(err) {
		return false, "", err
	}

	return true, container_id, nil
}

// Check if the input folder is mounted into the container
func IsContainerFolderMounted(folder string) *FindmntOutput {
	var findmnt_output bytes.Buffer // Initialize the output command buffer

	// Create the findmnt command, links the output and run
	findmnt := exec.Command("findmnt", "--target", folder, "--json")
	findmnt.Stdout = &findmnt_output
	err := findmnt.Run()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Unmarshal the output from findmnt into the structure
	output := new(FindmntOutput) // Create the output structure
	err = json.Unmarshal(findmnt_output.Bytes(), output)
	if err != nil {
		fmt.Printf("Error parsing the JSON: %s\n", err.Error())
		return nil
	}

	// Check if the target is the same as the input folder. This is a
	// necessary check, since findmnt with current folder returns always
	// something. For folders that are not bind mounted for example
	// there would be something referring to the root folder.
	if len(output.Filesystems) < 1 || output.Filesystems[0].Target != folder {
		return nil
	}

	return output
}

// Performs all the containers checks and returns the target folder
func PerformContainerChecks(cwd string) (string, string, string) {
	fmt.Println("------------------------- CONTAINER CHECKS ---------------------------------")

	// First we need to check if the current running environment is a docker
	// container or not. If it is then more diagnostic is necessary, otherwise
	// we can just returns the current working folder
	fmt.Print("[*] Is the current environment under a container? ")
	is_container, container_id, _ := IsContainerEnvironment()
	target_folder := cwd
	source_folder := cwd
	entry_path := cwd

	if !is_container {
		fmt.Println("No")
	} else {
		fmt.Printf("Yes.\nDETECTED CONTAINER ID: %s\n", container_id)

		// Check if the current working folder is mounted inside the container.
		// This checks gives also a mapping from the current working folder
		// to the host working folder, which is likely to be written into the
		// .git file if the target is inside git worktree
		fmt.Print("\n[*] Checking if the CWD is bind mounted: ")
		fs_output := IsContainerFolderMounted(cwd)

		var fs_data *Filesystem = nil
		if fs_output == nil {
			fmt.Println("No")
		} else {
			fs_data = &fs_output.Filesystems[0] // Select the only entry
			fmt.Println("Yes")
			fmt.Printf("   HOST SOURCE: %s\n", fs_data.Source)
			fmt.Printf("   CONTAINER TARGET: %s\n", fs_data.Target)
			parts := strings.Split(fs_data.Source, "[")
			source_folder = strings.TrimSuffix(parts[1], "]")
			entry_path = fs_data.Target
		}

		// Get the ccommits target folder from the environment variable.
		// According to the ccommits documentation, user can specify
		// the target folder using an env variable named CCOMMITS_WD
		fmt.Print("\n[*] Is the CCOMMITS_WD env variable set: ")
		target_folder, _ = GetContainerEnvironmentVariable("CCOMMITS_WD")
		if len(target_folder) < 1 {
			// If the environment variable is not set, then ask the user
			fmt.Println("No")
			fmt.Print("[*] Enter the (relative) target folder (Leave blank for CWD): ")
			fmt.Scanln(&target_folder)
			if len(target_folder) < 1 {
				target_folder = cwd // Set to the current working folder
			} else {
				target_folder = filepath.Join(cwd, target_folder) // Construct the path
			}
		} else {
			target_folder = filepath.Join(cwd, target_folder)
		}

		fmt.Printf("SELECTED TARGET FOLDER: %s\n", target_folder)
	}

	fmt.Println("----------------------------------------------------------------------------")

	return target_folder, source_folder, entry_path
}
